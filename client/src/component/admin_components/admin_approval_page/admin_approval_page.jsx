import React, { useState, useEffect } from "react";
import "./admin_approval_page.css";

const AdminApprovalPage = () => {
    const [registrations, setRegistrations] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState("");
    const [rejectReason, setRejectReason] = useState(""); // Состояние для причины отклонения
    const [selectedRegistrationId, setSelectedRegistrationId] = useState(null); // ID выбранной заявки
    const [selectedUser, setSelectedUser] = useState(null); // Выбранный пользователь

    useEffect(() => {
        fetch("http://localhost:5000/admin/course-registrations", {
            credentials: "include",
        })
            .then((response) => response.json())
            .then((data) => {
                console.log("Data from backend:", data); // Проверьте структуру данных
                if (data.error) {
                    setError(data.error);
                } else {
                    setRegistrations(data);
                }
                setIsLoading(false);
            })
            .catch((error) => {
                console.error("Error fetching registrations:", error);
                setError("Ошибка при загрузке заявок");
                setIsLoading(false);
            });
    }, []);

    const handleApprove = (registrationId) => {
        fetch(`http://localhost:5000/admin/approve-registration/${registrationId}`, {
            method: "POST",
            credentials: "include",
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.success) {
                    setRegistrations(
                        registrations.map((reg) =>
                            reg._id === registrationId
                                ? { ...reg, status: "Одобренный" }
                                : reg
                        )
                    );
                } else {
                    setError(data.message || "Ошибка при одобрении заявки");
                }
            })
            .catch((error) => {
                console.error("Error approving registration:", error);
                setError("Ошибка при одобрении заявки");
            });
    };

    const handleReject = (registrationId) => {
        setSelectedRegistrationId(registrationId);
    };

    const confirmReject = () => {
        if (!rejectReason) {
            setError("Укажите причину отклонения");
            return;
        }
        fetch(`http://localhost:5000/admin/reject-registration/${selectedRegistrationId}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            credentials: "include",
            body: JSON.stringify({ reason: rejectReason }), // Отправляем причину отклонения
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.success) {
                    setRegistrations(
                        registrations.map((reg) =>
                            reg._id === selectedRegistrationId
                                ? { ...reg, status: "Отклоненный", rejectReason: rejectReason }
                                : reg
                        )
                    );
                    setRejectReason(""); // Очищаем поле причины
                    setSelectedRegistrationId(null); // Закрываем модальное окно
                } else {
                    setError(data.message || "Ошибка при отклонении заявки");
                }
            })
            .catch((error) => {
                console.error("Error rejecting registration:", error);
                setError("Ошибка при отклонении заявки");
            });
    };

    const handleViewUser = (user) => {
        setSelectedUser(user);
    };

    const handleCloseUserInfo = () => {
        setSelectedUser(null);
    };

    if (isLoading) {
        return <div>Загрузка...</div>;
    }

    if (error) {
        return <div>{error}</div>;
    }

    return (
        <div className="admin-approval-page">
            <h1>Заявки на курсы</h1>
            <table>
                <thead>
                    <tr>
                        <th>Курс</th>
                        <th>Пользователь</th>
                        <th>Статус</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    {registrations.map((registration) => (
                        <tr key={registration._id}>
                            <td>{registration.courseTitle}</td>
                            <td>
                            <button onClick={() => handleViewUser({
                                    username: registration.userName,
                                    email: registration.userEmail,
                                    fullname: registration.userFullname,
                                    birthdate: registration.userBirthdate,
                                    birthplace: registration.userBirthplace,
                                    education: registration.userEducation,
                                    workplace: registration.userWorkplace,
                                    jobtitle: registration.userJobtitle,
                                    homeaddress: registration.userHomeaddress,
                                    phone: registration.userPhone,
                                    passportdata: registration.userPassportdata,
                                    snils: registration.userSnils,
                                })}>
                                    {registration.userName}
                                </button>
                            </td>
                            <td>{registration.status}</td>
                            <td>
                                {registration.status === "Ожидание" && (
                                    <>
                                        <button className="approve-btn" onClick={() => handleApprove(registration._id)}>
                                            Одобрить
                                        </button>
                                        <button className="reject-btn" onClick={() => handleReject(registration._id)}>
                                            Отклонить
                                        </button>
                                    </>
                                )}
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>

            {/* Модальное окно для ввода причины отклонения */}
            {selectedRegistrationId && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Укажите причину отклонения</h2>
                        <textarea
                            value={rejectReason}
                            onChange={(e) => setRejectReason(e.target.value)}
                            placeholder="Причина отклонения"
                        />
                        <button onClick={confirmReject}>Подтвердить</button>
                        <button onClick={() => setSelectedRegistrationId(null)}>Отмена</button>
                    </div>
                </div>
            )}

            {/* Модальное окно для просмотра информации о пользователе */}
            {selectedUser && (
                <div className="user-info-modal">
                    <div className="user-info-container">
                        <h3>Информация о пользователе</h3>
                        <table className="user-info-table">
                            <tbody>
                                <tr>
                                    <th>Логин</th>
                                    <td>{selectedUser.username || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Email</th>
                                    <td>{selectedUser.email || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>ФИО</th>
                                    <td>{selectedUser.fullname || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Дата рождения</th>
                                    <td>{selectedUser.birthdate || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Место рождения</th>
                                    <td>{selectedUser.birthplace || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Образование</th>
                                    <td>{selectedUser.education || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Место работы</th>
                                    <td>{selectedUser.workplace || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Должность</th>
                                    <td>{selectedUser.jobtitle || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Домашний адрес</th>
                                    <td>{selectedUser.homeaddress || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Номер телефона</th>
                                    <td>{selectedUser.phone || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Данные паспорта</th>
                                    <td>{selectedUser.passportdata || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>СНИЛС</th>
                                    <td>{selectedUser.snils || "Данные отсутствуют"}</td>
                                </tr>
                            </tbody>
                        </table>
                        <button onClick={handleCloseUserInfo}>Закрыть</button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AdminApprovalPage;