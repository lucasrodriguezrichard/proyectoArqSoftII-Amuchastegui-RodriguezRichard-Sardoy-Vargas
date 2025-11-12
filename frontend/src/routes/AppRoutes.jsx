import { Outlet, Route, Routes, Navigate } from 'react-router-dom';

import Home from '../pages/Home';
import Login from '../pages/Login';
import Register from '../pages/Register';
import ReservationDetails from '../pages/ReservationDetails';
import MyReservations from '../pages/MyReservations';
import Admin from '../pages/Admin';
import PrivateRoute from './PrivateRoute';
import AdminRoute from './AdminRoute';
import { Navbar } from '../components/common/Navbar';
import { Footer } from '../components/common/Footer';

const AppLayout = () => (
  <div className="flex min-h-screen flex-col bg-slate-50">
    <Navbar />
    <main className="flex-1">
      <Outlet />
    </main>
    <Footer />
  </div>
);

const AppRoutes = () => (
  <Routes>
    <Route element={<AppLayout />}>
      <Route path="/" element={<Home />} />
      <Route path="/reservations/:id" element={<ReservationDetails />} />
      <Route
        path="/my-reservations"
        element={
          <PrivateRoute>
            <MyReservations />
          </PrivateRoute>
        }
      />
      <Route
        path="/admin"
        element={
          <AdminRoute>
            <Admin />
          </AdminRoute>
        }
      />
    </Route>
    <Route path="/login" element={<Login />} />
    <Route path="/register" element={<Register />} />
    <Route path="*" element={<Navigate to="/" replace />} />
  </Routes>
);

export default AppRoutes;
