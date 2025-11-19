import { Request, Response } from 'express';
import User from '../models/User.js';
import jwt from 'jsonwebtoken';
import TryCatch from '../utils/TryCatch.js';
import { AuthenticatedRequest } from '../middleware/isAuth.js';

export const loginUser = TryCatch(async (req , res) => {
    const {name , email, image} = req.body;

    let user  =  await User.findOne({email});

    if(!user) {
        user  = await User.create({name, email, image});
    }

    const token = jwt.sign({user}, process.env.JWT_SECRET as string, {expiresIn: '5d'});

    res.status(200).json({message:"Login successful", token, user});
})

export const myProfile = TryCatch(async (req : AuthenticatedRequest, res) => {
    const user = req.user;
    res.json({user});
});