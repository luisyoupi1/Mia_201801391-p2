import PropTypes from 'prop-types'; // Primero, asegúrate de instalar prop-types con npm o yarn
import SessionContext from './SessionContext';
import { useState } from 'react';

export const SessionProvider = ({ children }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    const login = () => {
        setIsAuthenticated(true);
        sessionStorage.setItem('session', 'true');
    };

    const logout = () => {
        setIsAuthenticated(false);
        sessionStorage.removeItem('session');
    };

    return (
        <SessionContext.Provider value={{ isAuthenticated, login, logout }}>
            {children}
        </SessionContext.Provider>
    );
};

SessionProvider.propTypes = {
    children: PropTypes.node.isRequired // 'node' cubre cualquier cosa que pueda ser renderizada: números, strings, elementos o fragmentos.
};
