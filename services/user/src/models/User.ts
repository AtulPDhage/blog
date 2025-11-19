import mongoose, {Document, Schema} from "mongoose";

export interface IUser extends Document {
    name : string;
    email : string;
    image : string;
    instagram : string;
    linkedin : string;
    facebook : string;
    bio : string;
}

const schema : Schema<IUser> = new mongoose.Schema({
    name: {type: String, required: true},
    email: {type: String, required: true, unique: true},
    image: {type: String, required: true},
    instagram: {type: String, required: false},
    linkedin: {type: String, required: false},
    facebook: {type: String, required: false},
    bio: {type: String, required: false},
}, {timestamps: true});

const User = mongoose.model<IUser>("User", schema);

export default User;