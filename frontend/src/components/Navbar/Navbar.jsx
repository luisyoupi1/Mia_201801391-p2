import { FaCode } from "react-icons/fa";
import { LuFiles } from "react-icons/lu";
import { TbReportSearch } from "react-icons/tb";
import { BiSolidLogOut } from "react-icons/bi";
import { Link } from "react-router-dom";
import logo from "../../assets/consola.png"
import { useSession } from "../../session/useSession";



function Navbar() {

    const { isAuthenticated, logout } = useSession()

    return (
        <nav className="navbar navbar-expand-lg navbar-dark bg-dark py-4">
            <div className="container">
                <Link className="navbar-brand align-content-start" to={'/'}><img src={logo} alt="consola" className="w-logo" />MyTerminal</Link>
                <div className="align-content-end" id="navbarNav">
                    <ul className="navbar-nav fs-5">
                        <li className="nav-item">
                            <Link className="nav-link" to={'/'}>Comandos <FaCode /></Link>
                        </li>
                        <li className="nav-item">
                            <Link className="nav-link" to={'/files'}>Archivos <LuFiles /></Link>
                        </li>
                        <li className="nav-item">
                            <Link className="nav-link" to={'/reports'}>Reportes <TbReportSearch /></Link>
                        </li>
                        {
                            isAuthenticated ? (
                                <li className="nav-item">
                                    <button className="nav-link" onClick={logout}>Logout<BiSolidLogOut /></button>
                                </li>
                            ) : (
                                null
                            )
                        }

                    </ul>
                </div>
            </div>
        </nav>
    );
}

export default Navbar;