import express from 'express';
import { getUserProfile, loginUser, myProfile, updateUser } from '../controllers/user.js';
import { isAuth } from '../middleware/isAuth.js';

const router = express.Router();

router.post('/login', loginUser);
router.get('/me', isAuth, myProfile);
router.get('/user/:id', getUserProfile);
router.post('/user/update', isAuth, updateUser);

export default router;