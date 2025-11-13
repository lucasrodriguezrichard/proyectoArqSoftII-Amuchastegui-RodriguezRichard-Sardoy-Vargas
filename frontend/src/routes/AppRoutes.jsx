import { Outlet, Route, Routes, Navigate } from 'react-router-dom';

import Home from '../pages/Home';
import Login from '../pages/Login';
import Register from '../pages/Register';
import ReservationDetails from '../pages/ReservationDetails';
import MyReservations from '../pages/MyReservations';
import CreateReservation from '../pages/CreateReservation';
import Admin from '../pages/Admin';
import PrivateRoute from './PrivateRoute';
import AdminRoute from './AdminRoute';
import { Navbar } from '../components/common/Navbar';
import { Footer } from '../components/common/Footer';

const AppLayout = () => (
  <div className="relative min-h-screen bg-slate-50 dark:bg-slate-950">
    <div className="pointer-events-none absolute inset-0 opacity-80">
      <div className="absolute -top-32 left-1/2 h-64 w-64 -translate-x-1/2 rounded-full bg-primary-400/20 blur-3xl dark:bg-primary-700/30" />
      <div className="absolute bottom-0 right-0 h-72 w-72 translate-x-1/2 translate-y-1/2 rounded-full bg-sky-300/30 blur-3xl dark:bg-sky-500/20" />
    </div>
    <div className="relative flex min-h-screen flex-col">
      <Navbar />
      <main className="flex-1">
        <Outlet />
      </main>
      <Footer />
    </div>
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
        path="/create-reservation"
        element={
          <PrivateRoute>
            <CreateReservation />
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
