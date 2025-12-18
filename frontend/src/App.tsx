import { lazy } from 'react';
import { Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';

const Dashboard = lazy(() => import('./pages/Dashboard'));
const Traffic = lazy(() => import('./pages/Traffic'));
const ProxyManager = lazy(() => import('./pages/ProxyManager'));
const Rules = lazy(() => import('./pages/Rules'));
const Reports = lazy(() => import('./pages/Reports'));
const Settings = lazy(() => import('./pages/Settings'));
const OpenSource = lazy(() => import('./pages/OpenSource'));
const About = lazy(() => import('./pages/About'));
import { EventsOn } from '../wailsjs/runtime/runtime';
import { useNavigate } from 'react-router-dom';
import { useEffect } from 'react';

function App() {
    const navigate = useNavigate();

    useEffect(() => {
        const unsubscribe = EventsOn("navigate-to", (path: string) => {
            navigate(path);
        });

        const handleKeyDown = (e: KeyboardEvent) => {
            // Disable Ctrl+A or Cmd+A (Select All)
            if ((e.ctrlKey || e.metaKey) && e.key === 'a') {
                const target = e.target as HTMLElement;
                // Allow select all in inputs or textareas
                if (target.tagName !== 'INPUT' && target.tagName !== 'TEXTAREA' && !target.isContentEditable) {
                    e.preventDefault();
                }
            }
        };

        window.addEventListener('keydown', handleKeyDown);

        return () => {
            window.removeEventListener('keydown', handleKeyDown);
            // Unsubscribe if runtime supports it, otherwise it's fine for singleton app
        };
    }, [navigate]);

    return (
        <Routes>
            <Route path="/" element={<Layout />}>
                <Route index element={<Dashboard />} />
                <Route path="traffic" element={<Traffic />} />
                <Route path="proxy" element={<ProxyManager />} />
                <Route path="rules" element={<Rules />} />
                <Route path="reports" element={<Reports />} />
                <Route path="settings" element={<Settings />} />
                <Route path="opensource" element={<OpenSource />} />
                <Route path="about" element={<About />} />
            </Route>
        </Routes>
    )
}

export default App
