import { Link } from "react-router-dom";
import './header_top.component.css';

const Header_top = () => {
    return (
        <div className="header_top ">
            <nav>
                <ul>
                    <li><Link to="/document" className="a">Нормативные документы</Link></li>
                    <li><Link to="/teacher" className="a">Преподаватели</Link></li>
                    <li><Link to="/sveden" className="a">Сведения об образовательной организации</Link></li>
                    {/* <li><Link to="/" className="a">Контакты</Link></li>  */}
                </ul>
            </nav>
        </div>
    )

};

export default Header_top;