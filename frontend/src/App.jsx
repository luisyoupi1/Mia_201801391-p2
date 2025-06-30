import Navbar from "./components/Navbar/Navbar"
import Footer from "./components/Footer/Footer"
import Comander from "./components/Comander/Comander"
import File from "./components/File/File"
import Partition from "./components/Partition/Partition"
import Login from "./components/Login/Login"
import Explorer from "./components/Explorer/Explorer"
import Report from "./components/Report/Report"

import { Routes, Route, HashRouter } from "react-router-dom"
import { SessionProvider } from "./session/Session"

export const ENDPOINT = "http://3.147.32.143:4000";
sessionStorage.setItem('session', false);


function App() {

  return (
    <HashRouter>
      <SessionProvider>
        <Navbar />
        <Routes>
          <Route exact path="/" element={<Comander />} />
          <Route exact path="/files" element={<File />} />
          {/* Para usar queryParams colocamos (:nombreElemento) */}
          <Route exact path="/files/:driveletter" element={<Partition />} />
          <Route exact path="/files/:driveletter/:partition" element={<Login />} />
          <Route exact path="/files/:driveletter/:partition/explorer" element={<Explorer />} />
          <Route exact path="/reports" element={<Report />} />
        </Routes>
        <Footer />
      </SessionProvider>
    </HashRouter>
  )
}

export default App
