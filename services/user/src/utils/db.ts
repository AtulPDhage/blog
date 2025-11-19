import mongoose from "mongoose";

const connectDB = async () => {
    try {
        await mongoose.connect(process.env.MONGO_URI as string, {
            dbName : process.env.DB_NAME as string  ?? "MasterDB",
        });
        console.log("MongoDB connected successfully");
    }       
    catch (error) {
        console.error("MongoDB connection error:", error);
        process.exit(1);
    }
};

export default connectDB;