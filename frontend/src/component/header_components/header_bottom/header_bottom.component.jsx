import { Link } from "react-router-dom";
import './header_bottom.component.css';

const Header_bottom = () => {
    return (
        <div className="header_bottom">
            <nav>
                <ul>
                    <li><Link to="/professional" className="a">Профессиональная переподготовка</Link></li>
                    <li><Link to="/promotion" className="a">Повышение квалификации  </Link></li>
                    <li><Link to="/admin" className="a">Вход системы</Link></li>
                </ul>
            </nav>
        </div>
    )

};


export default Header_bottom;