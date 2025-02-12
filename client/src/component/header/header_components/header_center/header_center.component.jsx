import logo from '../../../../assets/logo_apt.png';
import address from '../../../../assets/address.svg';
import phone from '../../../../assets/phone.svg';
import email from '../../../../assets/email.svg';
import { Link } from 'react-router-dom';
import './header_ center.component.css';
const Header_center = () => {
    return (
        <div className="header_center">
            <div className="header_logo">
                <div className="logo"><Link to="/" className="a"><img src={logo} alt="" /></Link></div>
            </div>
            <Link to='/' className="header_info">
                <p className='black'>
                    МИНИСТЕРСТВО ОБРАЗОВАНИЯ И НАУКИ РЕСПУБЛИКИ ТАТАРСТАН
                </p>
                <p className="small black">
                    Государственное автономное профессиональное образовательное учреждение
                </p>
                <p className="apt_color bold">
                    «АЛЬМЕТЬЕВСКИЙ ПОЛИТЕХНИЧЕСКИЙ ТЕХНИКУМ»
                </p>
            </Link>
            <div className="header_addres">
                <div className="header_dirrector">
                    <p>
                        Директор: <span className="apt_color">И.Р.Багманов</span>
                    </p>
                </div>
                <div className="header_icons">
                    <img src={address} alt="" />
                    <p>
                        г.Альметьевск, ул.Мира, д.10
                    </p>
                </div>
                <div className="header_icons">
                    <img src={phone} alt="" />
                    <p>
                        8 (8553) 33-48-56, 39-99-19
                    </p>
                </div>
                <div className="header_icons">
                    <img src={email} alt="" />
                    <p>
                        info@almetpt.ru
                    </p>
                </div>

            </div>
        </div>
    )

};

export default Header_center;