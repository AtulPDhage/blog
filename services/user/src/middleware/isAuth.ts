import { Request, Response, NextFunction } from "express";
import jwt, { JwtPayload } from "jsonwebtoken";
import { IUser } from "../models/User.js";

export interface AuthenticatedRequest extends Request {
    user?: IUser | null;
}


 export const isAuth = async (req : AuthenticatedRequest, res : Response, next : NextFunction)
: Promise<void> => {
    try {
        const authHeader = req.headers.authorization;
        if (!authHeader || !authHeader.startsWith("Bearer ")) {
            res.status(401).json({ message: "Unauthorized" });
            return;
        }

        const token = authHeader.split(" ")[1];
        const decodedValue = jwt.verify(token, process.env.JWT_SECRET as string) as JwtPayload;

        if (!decodedValue || !decodedValue.user) {
            res.status(401).json({ message: "Unauthorized" });
            return;
        }

        req.user = decodedValue.user;
        next();

    }
    catch (error) {
        console.error("Authentication error:", error);
        res.status(401).json({ message: "Unauthorized" });
    }
}