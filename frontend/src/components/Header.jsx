import styles from "./styles/HomePage.module.css";
import React from "react";
import { Link } from 'react-router-dom';

const Header = () => {
    return (
        <header className={styles.header}>
            <h1 className={styles.title}>Salafi Fatawa App</h1>
            <nav className={styles.nav}>
                <ul>
                    <li>
                        <a href="/">Home</a>
                    </li>
                    <li>
                        <a href="/fatawa">Fatawa</a>
                    </li>
                    <li>
                        <a href="/login">Log In</a>
                    </li>
                    <li>
                        <a href="/signup">Sign Up</a>
                    </li>
                </ul>
            </nav>
        </header>
    );
};

export default Header;