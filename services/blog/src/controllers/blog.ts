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