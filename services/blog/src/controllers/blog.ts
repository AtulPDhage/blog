import { AuthenticatedRequest } from "../middleware/isAuth.js";
import { redisClient } from "../server.js";
import { sql } from "../utils/db.js";
import TryCatch from "../utils/TryCatch.js";
import axios from "axios";

export const getAllBlogs = TryCatch(async (req, res) => {
    const {searchQuery = "", category = ""} = req.query;

    const cacheKey = `blogs:${searchQuery}:${category}`;

    const cached =  await redisClient.get(cacheKey);

    if(cached){
        console.log("Serving from Redis cache");
        return res.status(200).json({ blogs: JSON.parse(cached) });
    }

    let blogs; 

    if(searchQuery && category){
        blogs  = await sql`SELECT * FROM blogs WHERE (title ILIKE ${'%' + searchQuery + '%'} OR description ILIKE ${'%' + searchQuery + '%'} ) AND category=${category} ORDER BY created_at DESC;`;
    }
    else if(searchQuery){
        blogs  = await sql`SELECT * FROM blogs WHERE (title ILIKE ${'%' + searchQuery + '%'} OR description ILIKE ${'%' + searchQuery + '%'} ) ORDER BY created_at DESC;`;
    }
    else if(category){
        blogs  = await sql`SELECT * FROM blogs  WHERE category=${category} ORDER BY created_at DESC;`; 
    }
    else{
        blogs  = await sql`SELECT * FROM blogs ORDER BY created_at DESC;`; 
    }
    console.log("Serving from Database");

    await redisClient.set(cacheKey, JSON.stringify(blogs), { EX: 3600 }); // Cache for 60 minutes 

    res.status(200).json({ blogs });
});

export const getSingleBlog = TryCatch(async (req, res) => {
    const { id } = req.params;

    const cacheKey = `blog:${id}`;

    const cached =  await redisClient.get(cacheKey);

    if(cached){
        console.log("Serving single blog from Redis cache");
        return res.status(200).json({ blog: JSON.parse(cached) });
    }

    const blog = await sql`SELECT * FROM blogs WHERE id = ${id};`;
    if (blog.length === 0) {
        return res.status(404).json({ message: "Blog not found" });
    }

    const {data}  = await axios.get(`${process.env.USER_SERVICE}/api/v1/user/${blog[0].author}`);

    const responseData =  {blog: blog[0] , author: data};
    await redisClient.set(cacheKey, JSON.stringify(responseData), { EX: 3600 }); // Cache for 60 minutes

    console.log("Serving single blog from Database");
    
    res.status(200).json(responseData);
});

export const addComment = TryCatch(async (req : AuthenticatedRequest, res) => {
    const {id : blogid} = req.params;
    const {comment}  = req.body;

    await sql`INSERT INTO comments(comment, blogid, userid, username) VALUES (${comment}, ${blogid}, ${req.user?._id}, ${req.user?.name}) RETURNING *`;
    res.json({
        message : "Comment Added",
    })
})

export const getAllComments = TryCatch(async (req , res) => {
    const {id} =  req.params;

    let comments = await sql`SELECT * FROM comments WHERE blogid = ${id} ORDER BY created_at DESC`

    return res.status(200).json(comments);
})

export const deleteComment  = TryCatch(async(req:AuthenticatedRequest, res) => {
    const {commentid} = req.params;
    const comment =  await sql`SELECT * FROM comments WHERE id = ${commentid}`;
    console.log(comment[0].userid, req.user?._id)
    if(comment[0].userid !== req.user?._id)
    {
        res.status(401).json({
            message: "You are not owner of this comment"
        });

        return;
    }

    await sql`DELETE FROM comments WHERE id = ${commentid}`;
    res.json({message:"comment deleted"});
})

export const saveBlog = TryCatch(async(req: AuthenticatedRequest, res)=> {
    const {blogid} = req.params;
    const userid = req.user?._id;

    if(!blogid || !userid){
        res.status(400).json({message : "Missing Blog id or user id"});
        return;
    }

    const existing = await sql`SELECT * FROM savedblogs WHERE userid = ${userid} AND blogid = ${blogid}`;
    if(existing.length ===  0){
        await sql`INSERT INTO savedblogs (userid, blogid) VALUES (${userid}, ${blogid})`;
        res.json({message : "Blog saved"});
        return;
    }else{
        await sql`DELETE FROM savedblogs WHERE blogid = ${blogid} AND userid = ${userid}`;
        res.json({message : "Blog removed from saved blogs"});
    }
})

export const getSavedBlogs =  TryCatch(async(req:AuthenticatedRequest, res)=>{
    const blogs = await sql`SELECT * FROM savedblogs WHERE userid = ${req.user?._id}`;
    res.json(blogs);
});