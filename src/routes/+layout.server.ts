import type { LayoutServerLoad } from './$types'
import { error } from '@sveltejs/kit'
import { connectMongodb } from '~/lib/server/db'

export const load: LayoutServerLoad = async ({ locals }) => {
  const connection = await connectMongodb()

  if (!connection) {
    throw error(500, 'Database connection failed')
  } else {
    console.log('MongoDB connection successful!')
  }

  return { theme: locals.theme, language: locals.language }
}
