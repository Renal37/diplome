import React from "react";
import "./sveden_component.css";
import bogmanov from "../../assets/bogmanov.jpeg";
import borodina from "../../assets/Borodina.jpeg";
import habarova from "../../assets/habarova.jpeg";
import petro from "../../assets/Petro.jpeg";
import petrov from "../../assets/Petrov.jpeg";
import agliulina from "../../assets/Аглиуллин.jpeg";
import vinocurova from "../../assets/Винокурова.jpeg";
import ganiev from "../../assets/Ганиев.jpeg";
import elistrova from "../../assets/Елистратова.jpeg";
import ramazanov from "../../assets/Рамазанов.jpeg";
import tazetdinova from "../../assets/Тазетдинова.jpeg";
import tihonova from "../../assets/Тиханова.jpeg";
import haidarova from "../../assets/Хайдарова.jpeg";

const Svede = () => {
    return (
        <div className="svede_container">
            {/* Основная информация о техникуме */}
            <div className="history-section">
                <p className="history-item">1964 год Основан как филиал Казанского строительного техникума, с целью обеспечения специалистами бурно развивающейся нефтяной и строительной индустрии Республики Татарстан.</p>
                <p className="history-item">1966 год Филиал преобразован в техникум газовой промышленности.</p>
                <p className="history-item">1972 год Техникум переименован в строительный, получив название – Альметьевский строительный техникум (АСТ).</p>
                <p className="history-item">1993 год Техникум получает статус политехнического, в связи с многопрофильностью учебного заведения.</p>
                <p className="history-item">май 1995 года Происходит объединение Альметьевского политехнического и Альметьевского нефтяного техникумов на базе первого.</p>
                <p className="history-item">15 января 2008 года Издано Постановление Кабинета Министров № 1 о преобразовании учебного заведения в автономное.</p>
            </div>

            {/* Учредители и руководство */}
            <div className="main_center_header">
                <div className="line"></div>
                <h1>Учредители и руководство</h1>
                <div className="line"></div>
            </div>
            <div className="staff-section">
                <div className="staff-card">
                    <img src={bogmanov} alt="" />
                    <strong>ДИРЕКТОР ТЕХНИКУМА</strong>
                    <p>Багманов Ильдар Раисович</p>
                </div>
                <div className="staff-card">
                    <img src={tazetdinova} alt="" />
                    <strong>Секретарь руководителя</strong>
                    <p>Тазетдинова Диана Решатовна</p>
                </div>
                <div className="staff-card">
                    <img src={ganiev} alt="" />
                    <strong>Первый заместитель директора по УВР</strong>
                    <p>Ганиев Рафаэль Маратович</p>
                </div>
                <div className="staff-card">
                    <img src={borodina} alt="" />
                    <strong>Заместитель директора по учебной работе</strong>
                    <p>Бородина Рамиля Мирзовна</p>
                </div>
                <div className="staff-card">
                    <img src={petrov} alt="" />
                    <strong>Заместитель директора по учебно-производственной работе</strong>
                    <p>Петров Евгений Григорьевич</p>
                </div>
                <div className="staff-card">
                    <img src={habarova} alt="" />
                    <strong>Заместитель директора по экономическим вопросам</strong>
                    <p>Хабарова Наталия Геннадьевна</p>
                </div>
                <div className="staff-card">
                    <img src={petro} alt="" />
                    <strong>Заместитель директора по административно-хозяйственной части</strong>
                    <p>Петров Радик Григорьевич</p>
                </div>
                <div className="staff-card">
                    <img src={ramazanov} alt="" />
                    <strong>Заместитель директора по информатизации</strong>
                    <p>Рамазанов Рустем Ринатович</p>
                </div>
                <div className="staff-card">
                    <img src={tihonova} alt="" />
                    <strong>Главный бухгалтер</strong>
                    <p>Тиханова Люция Равиловна</p>
                </div>
                <div className="staff-card">
                    <img src={vinocurova} alt="" />
                    <strong>Начальник отдела кадров</strong>
                    <p>Винокурова Елена Александровна</p>
                </div>
                <div className="staff-card">
                    <img src={elistrova} alt="" />
                    <strong>Юрисконсульт</strong>
                    <p>Елистратова Оксана Александровна</p>
                </div>
                <div className="staff-card">
                    <img src={haidarova} alt="" />
                    <strong>Начальник научно-методического отдела</strong>
                    <p>Хайдарова Назия Мисбаховна</p>
                </div>
                <div className="staff-card">
                    <img src={agliulina} alt="" />
                    <strong>Начальник ЦЦОД «IT-куб»</strong>
                    <p>Аглиуллин Ильнур Ульфатович</p>
                </div>
            </div>


            {/* Контакты */}
            <div className="main_center_header">
                <div className="line"></div>
                <h1>Контакты</h1>
                <div className="line"></div>
            </div>
            <div className="contact-section">
                <p>
                    <strong>Адрес:</strong> улица Мира 10, Альметьевск, Респ. Татарстан, 423457
                </p>
                <p>
                    <strong>Часы работы:</strong> понедельник–суббота: 07:00–20:00, воскресенье: Закрыто
                </p>
                <p>
                    <strong>Телефон:</strong> 8 (855) 333-48-56
                </p>
                <p>
                    <strong>Email:</strong> info@almetpt.ru
                </p>
            </div>
        </div>
    );
};

export default Svede;