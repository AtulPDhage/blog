import express from 'express';
import dotenv from 'dotenv';
import blogRoutes from './routes/blog.js';

dotenv.config();

const app = express();

app.use(express.json());

const PORT = process.env.PORT;

app.use('/api/v1', blogRoutes);

app.listen(PORT, () => {
    console.log(`Blog service is running on http://localhost:${PORT}`);
});