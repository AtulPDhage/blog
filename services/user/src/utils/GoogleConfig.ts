import {google} from 'googleapis';
import dotenv from 'dotenv';
dotenv.config();

const GOOGLE_CLIENT_ID = process.env.Google_Client_Id;
const GOOGLE_CLIENT_SECRET = process.env.Google_client_Secret;

export const oauth2Client = new google.auth.OAuth2(
  GOOGLE_CLIENT_ID,
  GOOGLE_CLIENT_SECRET,
  "postmessage"
);




