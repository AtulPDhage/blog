import { Request, Response } from 'express';
import User from '../models/User.js';
import jwt from 'jsonwebtoken';
import TryCatch from '../utils/TryCatch.js';
import { AuthenticatedRequest } from '../middleware/isAuth.js';
import { getBuffer } from '../utils/dataUri.js';
import { v2 as cloudinary } from 'cloudinary';
import { oauth2Client } from '../utils/GoogleConfig.js';
import axios from 'axios';

export const loginUser = TryCatch(async (req , res) => {
    
    const {code} = req.body;

    if(!code) {
        res.status(400).json({message: "Authorization code is required"});
        return;
    }

    const googleRes = await oauth2Client.getToken(code);

    oauth2Client.setCredentials(googleRes.tokens);

    const userRes = await axios.get(`https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=${googleRes.tokens.access_token}`);

    const {name , email, picture} = userRes.data;

    let user  =  await User.findOne({email});

    if(!user) {
        user  = await User.create({name, email, image : picture});
    }

    const token = jwt.sign({user}, process.env.JWT_SECRET as string, {expiresIn: '5d'});

    res.status(200).json({message:"Login successful", token, user});
})

export const myProfile = TryCatch(async (req : AuthenticatedRequest, res) => {
    const user = req.user;
    res.json({user});
});

export const getUserProfile =  TryCatch(async (req , res) => {
    const user  = await User.findById(req.params.id);
    console.log(user);
    if(!user) {
        res.status(404).json({message: "User not found"});
        return;
    }
    res.json({user});
});

export const updateUser = TryCatch(async (req : AuthenticatedRequest, res) => {
    const {name, instagram, linkedin, facebook, bio} =   req.body;
    const user  =  await User.findByIdAndUpdate(req.user?._id, {
        name, instagram, linkedin, facebook, bio
    }, {new: true});

     const token = jwt.sign({user}, process.env.JWT_SECRET as string, {expiresIn: '5d'});
    res.json({message: "Profile updated successfully", user, token});
});

export const updateProfilePicture = TryCatch(async (req : AuthenticatedRequest, res) => {
    const file  =  req.file;

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
        folder: "profile_pictures"
    });

    const user  =  await User.findByIdAndUpdate(req.user?._id, {
        image : cloud.secure_url
    }, {new: true}); 

    const token = jwt.sign({user}, process.env.JWT_SECRET as string, {expiresIn: '5d'});
    res.json({message: "Profile picture updated successfully", user, token});
});