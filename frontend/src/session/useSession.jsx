import { useContext } from 'react';
import SessionContext from './SessionContext';

export const useSession = () => useContext(SessionContext);
