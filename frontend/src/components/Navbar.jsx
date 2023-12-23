import { useState } from 'react'
import { NavLink } from 'react-router-dom'
import Hamburger from './Assets/menu.svg';
import './CSS/Navbar.css'



const Navbar = () => {
    const [showNavbar, setShowNavbar] = useState(false)

    const handleShowNavbar = () => {
        setShowNavbar(!showNavbar)
    }

    return (
        <nav className="navbar">
            <div className="container">
                <div className="logo">
                    <h1 className='text-3xl font-bold'>NekoNotes</h1>
                </div>
                <div className="menu-icon" onClick={handleShowNavbar}>
                    <img src={Hamburger} alt="Brand" />
                </div>
                <div className={`nav-elements  ${showNavbar && 'active'}`}>
                    <ul>
                        <li>
                            <NavLink to="/">Home</NavLink>
                        </li>
                        <li>
                            <NavLink to="/notes">Notes</NavLink>
                        </li>
                        <li>
                            <NavLink to="/login">Login</NavLink>
                        </li>
                    </ul>
                </div>
            </div>
        </nav>
    )
}

export default Navbar