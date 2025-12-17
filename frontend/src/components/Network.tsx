import { useTranslation } from 'react-i18next';

const Network = () => {
    const { t } = useTranslation();
    return (
        <div style={{ padding: '20px' }}>
            <h1>{t('network.title')}</h1>
            <p>{t('network.content')}</p>
        </div>
    );
};

export default Network;
