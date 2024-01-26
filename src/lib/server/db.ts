import { KUNGALGAME_MONGODB_URI } from '$env/static/private'
import mongoose from 'mongoose'

export const connectMongodb = async () => {
  try {
    return await mongoose.connect(KUNGALGAME_MONGODB_URI)
  } catch (error) {
    console.log(error)
  }
}
