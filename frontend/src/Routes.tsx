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
import { ExamplePage } from './components/pages/Example'

const AdminPage = lazy(() => import('@/components/pages/Admin'))

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<App />} errorElement={<ErrorPage />}>
      <Route path="/" element={<ExamplePage />} />
      <Route path="blog" element={<BlogListPage />} />
      <Route path="blog/:id" element={<BlogDetailPage />} />
      <Route path="about" element={<AboutPage />} />
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
