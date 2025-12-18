import { Loader2 } from 'lucide-react';

const Loading = () => {
    return (
        <div className="flex h-full w-full items-center justify-center p-8">
            <Loader2 className="animate-spin text-muted-foreground" size={32} />
        </div>
    );
};

export default Loading;
