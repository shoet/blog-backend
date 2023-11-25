import {
  Outlet,
  Route,
  RouterProvider,
  createBrowserRouter,
  createRoutesFromElements,
} from 'react-router-dom'
import { AboutPage } from './components/pages/About'
import { ErrorPage } from './components/pages/Error'
import App from './App'
import { BlogListPage } from './components/pages/BlogList'
import { BlogDetailPage } from './components/pages/BlogDetail'
import { SearchPage } from './components/pages/Search'
import { BlogEditPage } from './components/pages/BlogEdit'
import { BlogPostPage } from './components/pages/BlogPost'
import { VSplit } from './components/templates/VSplit'
import { SideContent } from './components/organisms/SideContent'
import { BaseLayout } from './components/templates/BaseLayout'
import { LoginPage } from './components/pages/Login'
import AdminPage from './components/pages/Admin'
import { PortfolioPage } from './components/pages/Portfolio'

const NormalLayout = () => {
  return (
    <BaseLayout>
      <VSplit MainContent={<Outlet />} SubContent={<SideContent />} />
    </BaseLayout>
  )
}

const AdminLayout = () => {
  return (
    <BaseLayout>
      <Outlet />
    </BaseLayout>
  )
}

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<App />} errorElement={<ErrorPage />}>
      <Route path="" element={<NormalLayout />}>
        <Route path="" element={<BlogListPage />} />
        <Route path=":id" element={<BlogDetailPage />} />
        <Route path="search" element={<SearchPage />} />
        <Route path="about" element={<AboutPage />} />
        <Route path="portfolio" element={<PortfolioPage />} />
        <Route path="new" element={<BlogPostPage />} />
        <Route path=":id/edit" element={<BlogEditPage />} />
      </Route>
      <Route path="admin" element={<AdminLayout />}>
        <Route path="" element={<AdminPage />} />
        <Route path="login" element={<LoginPage />} />
      </Route>
    </Route>,
  ),
)

export const Routes = () => <RouterProvider router={router} />
