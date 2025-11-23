import { AuthenticatedRequest } from "../middlewares/isAuth.js";
import { getBuffer } from "../utils/dataUri.js";
import { sql } from "../utils/db.js";
import TryCatch from "../utils/TryCatch.js";
import { v2 as cloudinary } from "cloudinary";

export const createBlog = TryCatch(async (req : AuthenticatedRequest, res) => {
    const {title, description, blogcontent, category} = req.body;
    const file  = req.file;

    if(!file) {
        res.status(400).json({message: "No file uploaded"});
        return;
    }

    const fileBuffer =  getBuffer(file);
    if(!fileBuffer || !fileBuffer.content) {
        res.status(500).json({message: "Error processing file"});
        return;
    }

    const cloud   = await cloudinary.uploader.upload(fileBuffer.content, {
        folder: "blogs"
    });

    const result =  await sql`INSERT INTO blogs (title, description, image, blogcontent, category, author) VALUES (${title}, ${description}, ${cloud.secure_url}, ${blogcontent}, ${category}, ${req.user?._id}) RETURNING *`;

    res.status(201).json({message: "Blog created successfully", blog: result[0]});
}); 

export const updateBlog = TryCatch(async (req : AuthenticatedRequest, res) => {
    const {id} = req.params;
    const {title, description, blogcontent, category} = req.body;
    const file  = req.file; 

    const blog  = await sql`SELECT * FROM blogs WHERE id = ${id}`;
    if(blog.length === 0) {
        res.status(404).json({message: "Blog not found"});
        return;
    }

    if(blog[0].author !== req.user?._id) {
        res.status(403).json({message: "You are not authorized to update this blog"});
        return;
    }
    let imageUrl = blog[0].image;

    if(file) {
        const fileBuffer =  getBuffer(file);
        if(!fileBuffer || !fileBuffer.content) {
            res.status(500).json({message: "Error processing file"});
            return;
        }
        const cloud   = await cloudinary.uploader.upload(fileBuffer.content, {
            folder: "blogs"
        });
        imageUrl = cloud.secure_url;
    }
    const result =  await sql`UPDATE blogs SET title = ${title || blog[0].title}, description = ${description || blog[0].description}, image = ${imageUrl}, blogcontent = ${blogcontent || blog[0].blogcontent}, category = ${category || blog[0].category} WHERE id = ${id} RETURNING *`;

    res.status(200).json({
        message: "Blog updated successfully",
        blog: result[0]
    });
});

export const deleteBlog = TryCatch(async (req : AuthenticatedRequest, res) => {
    const {id} = req.params;
    const blog  = await sql`SELECT * FROM blogs WHERE id = ${id}`;
    if(blog.length === 0) {
        res.status(404).json({message: "Blog not found"});
        return;
    }
    if(blog[0].author !== req.user?._id) {
        res.status(403).json({message: "You are not authorized to delete this blog"});
        return;
    }
    await sql`DELETE FROM blogs WHERE id = ${id}`;
    await sql`DELETE FROM comments WHERE blogid = ${id}`;
    await sql`DELETE FROM savedblogs WHERE blogid = ${id}`;

    res.status(200).json({message: "Blog deleted successfully"});
});