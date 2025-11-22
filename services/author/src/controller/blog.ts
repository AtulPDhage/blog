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