import { useState } from 'react';
import PropTypes from 'prop-types';
import SessionContext from './SessionContext';

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
    children: PropTypes.node.isRequired
};
