import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import "./admin_check_profile.css";

const AdminCheckProfile = () => {
    const [users, setUsers] = useState([]);
    const [fioSearchQuery, setFioSearchQuery] = useState("");
    const [loginSearchQuery, setLoginSearchQuery] = useState("");
    const [showFullInfo, setShowFullInfo] = useState(false);
    const location = useLocation();
    const navigate = useNavigate();

    useEffect(() => {
        const searchParams = new URLSearchParams(location.search);
        const username = searchParams.get("username");
        if (username) {
            setLoginSearchQuery(username);
        }

        fetch("http://localhost:5000/users")
            .then((res) => res.json())
            .then((data) => setUsers(data));
    }, [location.search]);

    const filteredUsers = users.filter((user) => {
        const fullName = [
            user.lastname || "",
            user.firstname || "",
            user.middleName || ""
        ].join(" ").toLowerCase();
        const login = user.username?.toLowerCase() || "";
        console.log(user.education);

        return (
            (!fioSearchQuery || fullName.includes(fioSearchQuery.toLowerCase())) &&
            (!loginSearchQuery || login.includes(loginSearchQuery.toLowerCase()))
        );
    });

    return (
        <div className="admin-delete-component">
            <div className="controls">
                <div className="search-container">
                    <input
                        type="text"
                        placeholder="Поиск по ФИО"
                        value={fioSearchQuery}
                        onChange={(e) => setFioSearchQuery(e.target.value)}
                        className="search-input"
                    />
                    <input
                        type="text"
                        placeholder="Поиск по логину"
                        value={loginSearchQuery}
                        onChange={(e) => setLoginSearchQuery(e.target.value)}
                        className="search-input"
                    />
                    <button
                        className="toggle-info-button"
                        onClick={() => setShowFullInfo(!showFullInfo)}
                    >
                        {showFullInfo ? "Скрыть детали" : "Показать детали"}
                    </button>
                </div>
            </div>
            <table className="course-table">
                <thead>
                    <tr>
                        <th>Логин</th>
                        <th>Email</th>
                        <th>ФИО</th>
                        {showFullInfo && (
                            <>
                                <th>Дата рождения</th>
                                <th>Место рождения</th>
                                <th>Образование</th>
                                <th>Место работы</th>
                                <th>Должность</th>
                                <th>Домашний адрес</th>
                                <th>Номер телефона</th>
                                <th>Данные паспорта</th>
                                <th>СНИЛС</th>
                                {/* <th>Согласие на обработку</th>
                                <th>Контракт загружен</th> */}
                            </>
                        )}
                    </tr>
                </thead>
                <tbody>
                    {filteredUsers.length > 0 ? (
                        filteredUsers.map((user) => (
                            <tr key={user._id}>
                                <td>{user.username || "Данные отсутствуют"}</td>
                                <td>{user.email || "Данные отсутствуют"}</td>
                                <td>
                                    {user.lastname && user.firstname && user.middlename
                                        ? `${user.lastname} ${user.firstname} ${user.middlename}`
                                        : "Данные отсутствуют"}
                                </td>
                                {showFullInfo && (
                                    <>
                                        <td>{user.birthdate || "Данные отсутствуют"}</td>
                                        <td>{user.birthplace || "Данные отсутствуют"}</td>
                                        <td>{user.education?.name || "Данные отсутствуют"}</td>
                                        <td>{user.workplace || "Данные отсутствуют"}</td>
                                        <td>{user.jobtitle || "Данные отсутствуют"}</td>
                                        <td>{user.homeaddress || "Данные отсутствуют"}</td>
                                        <td>{user.phone || "Данные отсутствуют"}</td>
                                        <td>{user.passportdata || "Данные отсутствуют"}</td>
                                        <td>{user.snils || "Данные отсутствуют"}</td>
                                        {/* <td>{user.agreetoprocessing ? "Да" : "Нет"}</td>
                                        <td>{user.contractUploaded ? "Да" : "Нет"}</td> */}
                                    </>
                                )}
                            </tr>
                        ))
                    ) : (
                        <tr>
                            <td colSpan={showFullInfo ? 14 : 3} className="no-data-message">
                                Нет данных для отображения
                            </td>
                        </tr>
                    )}
                </tbody>
            </table>
        </div>
    );
};

export default AdminCheckProfile;