import { Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import Traffic from './pages/Traffic';
import ProxyManager from './pages/ProxyManager';
import Rules from './pages/Rules';
import Reports from './pages/Reports';
import Settings from './pages/Settings';

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
