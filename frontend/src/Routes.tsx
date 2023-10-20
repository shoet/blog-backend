import {
  Route,
  RouterProvider,
  createBrowserRouter,
  createRoutesFromElements,
} from 'react-router-dom'
import { AboutPage } from './components/pages/About'
import { ErrorPage } from './components/pages/Error'
import App from './App'
import { Suspense, lazy } from 'react'
import { BlogListPage } from './components/pages/BlogList'
import { BlogDetailPage } from './components/pages/BlogDetail'
import { SearchPage } from './components/pages/Search'
import { BlogEditPage } from './components/pages/BlogEdit'
import { BlogPostPage } from './components/pages/BlogPost'

const AdminPage = lazy(() => import('@/components/pages/Admin'))

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<App />} errorElement={<ErrorPage />}>
      <Route path="" element={<BlogListPage />} />
      <Route path=":id" element={<BlogDetailPage />} />
      <Route path="search" element={<SearchPage />} />
      <Route path="about" element={<AboutPage />} />
      <Route path="new" element={<BlogPostPage />} />
      <Route path=":id/edit" element={<BlogEditPage />} />
      <Route
        path="admin"
        element={
          <Suspense fallback={<div>loading...</div>}>
            <AdminPage />
          </Suspense>
        }
      />
    </Route>,
  ),
)

export const Routes = () => <RouterProvider router={router} />
