import { Suspense } from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Loading from './common/Loading';

const Layout = () => {
    return (
        <div className="flex h-screen w-screen overflow-hidden bg-background text-foreground pointer-events-none">
            <Sidebar />
            <div className="flex-1 overflow-y-auto bg-background/50 relative pointer-events-auto">
                <Suspense fallback={<Loading />}>
                    <Outlet />
                </Suspense>
            </div>
        </div>
    );
};

export default Layout;
