import express from 'express';
import dotenv from 'dotenv';
import blogRoutes from './routes/blog.js';
import {createClient} from 'redis';

dotenv.config();

const app = express();

app.use(express.json());

const PORT = process.env.PORT;

export const redisClient = createClient({
    url: process.env.REDIS_URL ,
});

redisClient
  .connect()
  .then(() => {console.log('Connected to Redis');})
  .catch((err) => console.error("Redis connection error:", err));

app.use('/api/v1', blogRoutes);

app.listen(PORT, () => {
    console.log(`Blog service is running on http://localhost:${PORT}`);
});