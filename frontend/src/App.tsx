import { lazy } from 'react';
import { Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';

const Dashboard = lazy(() => import('./pages/Dashboard'));
const Traffic = lazy(() => import('./pages/Traffic'));
const ProxyManager = lazy(() => import('./pages/ProxyManager'));
const Rules = lazy(() => import('./pages/Rules'));
const Reports = lazy(() => import('./pages/Reports'));
const Settings = lazy(() => import('./pages/Settings'));

function App() {
    return (
        <Routes>
            <Route path="/" element={<Layout />}>
                <Route index element={<Dashboard />} />
                <Route path="traffic" element={<Traffic />} />
                <Route path="proxy" element={<ProxyManager />} />
                <Route path="rules" element={<Rules />} />
                <Route path="reports" element={<Reports />} />
                <Route path="settings" element={<Settings />} />
            </Route>
        </Routes>
    )
}

export default App
