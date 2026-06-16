import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import AuthLayout from './layouts/AuthLayout';
import MainLayout from './layouts/MainLayout';
import Login from './pages/Login';
import Register from './pages/Register';
import Feed from './pages/Feed';
import './index.css';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Auth Routes */}
        <Route element={<AuthLayout />}>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
        </Route>

        {/* Protected Routes */}
        <Route element={<MainLayout />}>
          <Route path="/feed" element={<Feed />} />
          <Route path="/groups" element={<div>Groups Page</div>} />
          <Route path="/messages" element={<div>Messages Page</div>} />
          <Route path="/notifications" element={<div>Notifications Page</div>} />
          <Route path="/profile" element={<div>Profile Page</div>} />
        </Route>

        {/* Redirects */}
        <Route path="/" element={<Navigate to="/feed" replace />} />
        <Route path="*" element={<Navigate to="/feed" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
