import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';

const Layout = () => {
    return (
        <div className="flex h-screen w-screen overflow-hidden bg-background text-foreground">
            <Sidebar />
            <div className="flex-1 overflow-y-auto bg-background/50 relative">
                <Outlet />
            </div>
        </div>
    );
};

export default Layout;
