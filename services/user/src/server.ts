import express from 'express';
import dotenv from 'dotenv';
import connectDB from './utils/db.js';
import userRoutes from './routes/user.js';
import { v2 as cloudinary } from 'cloudinary';
import cors from 'cors';

dotenv.config();

 // Configuration of cloudinary
    cloudinary.config({ 
        cloud_name: process.env.Cloud_Name, 
        api_key: process.env.Cloud_Api_Key, 
        api_secret: process.env.Cloud_Api_Secret
    });

const app = express();

app.use(express.json());
app.use(cors());

connectDB();

app.use('/api/v1', userRoutes);

const port  =  process.env.PORT;

app.listen(port, () => {
  console.log(`User service is running on port ${port} at http://localhost:${port}`);
});