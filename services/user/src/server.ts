import express from 'express';
import dotenv from 'dotenv';
import connectDB from './utils/db.js';
import userRoutes from './routes/user.js';

dotenv.config();

const app = express();

app.use(express.json());

connectDB();

app.use('/api/v1', userRoutes);

const port  =  process.env.PORT;

app.listen(port, () => {
  console.log(`User service is running on port ${port} at http://localhost:${port}`);
});