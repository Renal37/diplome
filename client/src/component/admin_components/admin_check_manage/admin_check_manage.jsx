import React, { useState, useEffect } from "react";
import "./admin_check_manage.css";

const AdminCoursesManagement = () => {
    const [registrations, setRegistrations] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState("");
    const [rejectReason, setRejectReason] = useState("");
    const [selectedRejectRegistrationId, setSelectedRejectRegistrationId] = useState(null);
    const [selectedPdfRegistrationId, setSelectedPdfRegistrationId] = useState(null);
    const [pdfUrl, setPdfUrl] = useState("");
    const [selectedUser, setSelectedUser] = useState(null);
    const [filterStatus, setFilterStatus] = useState("all");
    const [groups, setGroups] = useState([]);
    const [selectedRegistrations, setSelectedRegistrations] = useState([]);
    const [searchUserName, setSearchUserName] = useState("");
    const [searchCourseTitle, setSearchCourseTitle] = useState("");

    const deletableStatuses = ["Ожидание", "Одобренный", "Отклоненный"];

    const filteredRegistrations = registrations.filter(registration => {
        const matchesStatus = filterStatus === "all" || registration.status === filterStatus;
        const matchesUserName = registration.userName.toLowerCase().includes(searchUserName.toLowerCase());
        const matchesCourseTitle = registration.courseTitle.toLowerCase().includes(searchCourseTitle.toLowerCase());
        return matchesStatus && matchesUserName && matchesCourseTitle;
    });

    const handleMassApprove = async () => {
        if (selectedRegistrations.length === 0) {
            setError("Выберите хотя бы одну заявку");
            return;
        }

        const allPending = selectedRegistrations.every(id => {
            const registration = registrations.find(reg => reg._id === id);
            return registration?.status === "Ожидание";
        });

        if (!allPending) {
            setError("Одобрять можно только заявки 'Ожидание'");
            return;
        }

        try {
            const approvePromises = selectedRegistrations.map(id =>
                fetch(`http://localhost:5000/admin/approve-registration/${id}`, {
                    method: "POST",
                    credentials: "include",
                }).then(res => res.json())
            );

            const results = await Promise.all(approvePromises);
            const allSuccess = results.every(res => res.success);
            if (!allSuccess) {
                throw new Error("Ошибка при одобрении заявок");
            }

            setRegistrations(prev =>
                prev.map(reg =>
                    selectedRegistrations.includes(reg._id) && reg.status === "Ожидание"
                        ? { ...reg, status: "Одобренный" }
                        : reg
                )
            );
            setSelectedRegistrations([]);
        } catch (err) {
            console.error("Ошибка массового одобрения:", err);
            setError(err.message || "Ошибка при одобрении");
        }
    };

    const handleMassDelete = async () => {
        if (selectedRegistrations.length === 0) {
            setError("Выберите хотя бы одну заявку");
            return;
        }

        const allDeletable = selectedRegistrations.every(id => {
            const registration = registrations.find(reg => reg._id === id);
            return deletableStatuses.includes(registration?.status);
        });

        if (!allDeletable) {
            setError("Удалять можно только 'Ожидание', 'Одобренный', 'Отклоненный'");
            return;
        }

        try {
            const deletePromises = selectedRegistrations.map(id =>
                fetch(`http://localhost:5000/admin/delete-registration/${id}`, {
                    method: "DELETE",
                    credentials: "include",
                }).then(res => res.json())
            );

            const results = await Promise.all(deletePromises);
            const allSuccess = results.every(res => res.success);
            if (!allSuccess) {
                throw new Error("Ошибка при удалении заявок");
            }

            setRegistrations(prev => prev.filter(reg => !selectedRegistrations.includes(reg._id)));
            setSelectedRegistrations([]);
        } catch (err) {
            console.error("Ошибка массового удаления:", err);
            setError(err.message || "Ошибка при удалении");
        }
    };

    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const response = await fetch("http://localhost:5000/groups", { credentials: "include" });
                if (!response.ok) throw new Error(`Ошибка ${response.status}: Не удалось загрузить группы`);
                const data = await response.json();
                setGroups(data.groups || []);
            } catch (err) {
                console.error("Ошибка загрузки групп:", err);
                setError(err.message || "Ошибка при загрузке групп");
            }
        };

        const fetchRegistrations = async () => {
            try {
                const response = await fetch("http://localhost:5000/admin/course-registrations", { credentials: "include" });
                if (!response.ok) throw new Error(`Ошибка ${response.status}: Не удалось загрузить заявки`);
                const data = await response.json();
                if (data.error) {
                    setError(data.error);
                } else {
                    setRegistrations(data || []);
                }
                setIsLoading(false);
            } catch (err) {
                console.error("Ошибка загрузки заявок:", err);
                setError(err.message || "Ошибка при загрузке заявок");
                setIsLoading(false);
            }
        };

        fetchGroups();
        fetchRegistrations();
    }, []);

    useEffect(() => {
        const handleKeyDown = event => {
            if (event.key === "Escape") {
                setSelectedRejectRegistrationId(null);
                setSelectedPdfRegistrationId(null);
                setSelectedUser(null);
                setRejectReason("");
                setPdfUrl("");
            }
        };

        document.addEventListener("keydown", handleKeyDown);
        return () => document.removeEventListener("keydown", handleKeyDown);
    }, []);

    const handleApprove = async (registrationId) => {
        try {
            const response = await fetch(`http://localhost:5000/admin/approve-registration/${registrationId}`, {
                method: "POST",
                credentials: "include",
            });
            const data = await response.json();
            if (!data.success) {
                setError(data.error || "Ошибка при одобрении заявки");
                return;
            }
            setRegistrations(prev =>
                prev.map(reg =>
                    reg._id === registrationId ? { ...reg, status: "Одобренный" } : reg
                )
            );
        } catch (err) {
            console.error("Ошибка одобрения заявки:", err);
            setError(err.message || "Ошибка при одобрении заявки");
        }
    };

    const handleReject = (registrationId) => {
        setSelectedRejectRegistrationId(registrationId);
        setRejectReason("");
    };

    const confirmReject = async () => {
        if (!rejectReason) {
            setError("Укажите причину отклонения");
            return;
        }
        try {
            const response = await fetch(`http://localhost:5000/admin/reject-registration/${selectedRejectRegistrationId}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ reason: rejectReason })
            });
            const data = await response.json();
            if (!data.success) {
                setError(data.error || "Ошибка при отклонении заявки");
                return;
            }
            setRegistrations(prev =>
                prev.map(reg =>
                    reg._id === selectedRejectRegistrationId
                        ? { ...reg, status: "Отклоненный", rejectReason }
                        : reg
                )
            );
            setSelectedRejectRegistrationId(null);
            setRejectReason("");
        } catch (err) {
            console.error("Ошибка отклонения заявки:", err);
            setError(err.message || "Ошибка при отклонении заявки");
        }
    };

    const handleViewPdf = (registrationId) => {
        setSelectedPdfRegistrationId(registrationId);
        setPdfUrl(`http://localhost:5000/user/view-contract/${registrationId}`);
    };

    const handleApprovePdf = async (registrationId) => {
        try {
            const response = await fetch(`http://localhost:5000/admin/approve-contract/${registrationId}`, {
                method: "POST",
                credentials: "include",
            });
            const data = await response.json();
            if (!data.success) {
                setError(data.error || "Ошибка при принятии договора");
                return;
            }
            setRegistrations(prev =>
                prev.map(reg =>
                    reg._id === registrationId ? { ...reg, status: "Принят" } : reg
                )
            );
            setSelectedPdfRegistrationId(null);
            setPdfUrl("");
        } catch (err) {
            console.error("Ошибка принятия договора:", err);
            setError(err.message || "Ошибка при принятии договора");
        }
    };

    const handleAssignGroup = async (registrationId, groupId) => {
        if (!groupId) {
            setError("Выберите группу");
            return;
        }
        try {
            const response = await fetch(`http://localhost:5000/admin/assign-group/${registrationId}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ groupId })
            });
            const data = await response.json();
            if (!data.success) {
                setError(data.error || "Ошибка при привязке к группе");
                return;
            }
            setRegistrations(prev =>
                prev.map(reg =>
                    reg._id === registrationId ? { ...reg, groupId } : reg
                )
            );
        } catch (err) {
            console.error("Ошибка привязки к группе:", err);
            setError(err.message || "Ошибка при привязке к группе");
        }
    };

    const handleViewUser = (user) => {
        setSelectedUser(user);
    };

    const handleCloseUserInfo = () => {
        setSelectedUser(null);
    };

    const handleViewConsent = (userId) => {
        window.open(`http://localhost:5000/user/view-consent/${userId}`, '_blank');
    };

    const handleDelete = async (registrationId) => {
        const registration = registrations.find(reg => reg._id === registrationId);
        if (!deletableStatuses.includes(registration.status)) {
            setError("Удалять можно только 'Ожидание', 'Одобренный', 'Отклоненный'");
            return;
        }
        try {
            const response = await fetch(`http://localhost:5000/admin/delete-registration/${registrationId}`, {
                method: "DELETE",
                credentials: "include",
            });
            const data = await response.json();
            if (!data.success) {
                setError(data.error || "Ошибка при удалении");
                return;
            }
            setRegistrations(prev => prev.filter(reg => reg._id !== registrationId));
        } catch (err) {
            console.error("Ошибка удаления заявки:", err);
            setError(err.message || "Ошибка при удалении заявки");
        }
    };

    if (isLoading) return <div>Загрузка...</div>;

    return (
        <div className="admin-approval-page">
            {error && <p className="error-message">{error}</p>}
            <div className="controls">
                <div className="search-container">
                    <input
                        type="text"
                        value={searchUserName}
                        onChange={e => setSearchUserName(e.target.value)}
                        placeholder="Поиск по пользователю"
                        className="search-input"
                    />
                    <input
                        type="text"
                        value={searchCourseTitle}
                        onChange={e => setSearchCourseTitle(e.target.value)}
                        placeholder="Поиск по курсу"
                        className="search-input"
                    />
                    <div className="filter-section">
                        <label>Фильтр по статусу:</label>
                        <select value={filterStatus} onChange={e => setFilterStatus(e.target.value)}>
                            <option value="all">Все</option>
                            <option value="Ожидание">Ожидание</option>
                            <option value="Одобренный">Одобренный</option>
                            <option value="Отклоненный">Отклоненный</option>
                            <option value="Принят">Принят</option>
                            <option value="Оплаченный">Оплаченный</option>
                            <option value="Проходит курс">Проходит курс</option>
                            <option value="Отчисленный">Отчисленный</option>
                            <option value="Завершил">Завершил</option>
                        </select>
                    </div>
                </div>
            </div>
            <div className="mass-actions">
                <button
                    className="approve-btn"
                    onClick={handleMassApprove}
                    disabled={selectedRegistrations.length === 0 || !selectedRegistrations.every(id => {
                        const reg = registrations.find(r => r._id === id);
                        return reg.status === "Ожидание";
                    })}
                >
                    Одобрить выбранные
                </button>
                <button
                    className="reject-btn"
                    onClick={handleMassDelete}
                    disabled={selectedRegistrations.length === 0 || !selectedRegistrations.every(id => {
                        const reg = registrations.find(r => r._id === id);
                        return deletableStatuses.includes(reg.status);
                    })}
                >
                    Удалить выбранные
                </button>
            </div>
            <table>
                <thead>
                    <tr>
                        <th>Курс</th>
                        <th>Пользователь</th>
                        <th>Группа</th>
                        <th>Статус</th>
                        <th>Причина</th>
                        <th>Цена</th>
                        <th>Выборка</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    {filteredRegistrations.map(registration => {
                        const availableGroups = groups.filter(group => group.courseId === registration.courseId);

                        return (
                            <tr key={registration._id}>
                                <td className="course-title">{registration.courseTitle}</td>
                                <td>
                                    <button
                                        onClick={() =>
                                            handleViewUser({
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
                                            })}
                                    >
                                        {registration.userName}
                                    </button>
                                </td>
                                <td>
                                    {registration.status === "Оплаченный" ? (
                                        <select
                                            value={registration.groupId || ""}
                                            onChange={e => handleAssignGroup(registration._id, e.target.value)}
                                        >
                                            <option value="">Выберите группу</option>
                                            {availableGroups.map(group => (
                                                <option key={group._id} value={group._id}>
                                                    {group.groupName}
                                                </option>
                                            ))}
                                        </select>
                                    ) : (
                                        <span>{groups.find(g => g._id === registration.groupId)?.groupName || "Нет группы"}</span>
                                    )}
                                </td>
                                <td>{registration.status}</td>
                                <td>
                                    {registration.rejectReason && (
                                        <div className="reason-text">{registration.rejectReason}</div>
                                    )}
                                </td>
                                <td>{registration.price ? `${registration.price} руб.` : "Не указана"}</td>
                                <td>
                                    <input
                                        type="checkbox"
                                        checked={selectedRegistrations.includes(registration._id)}
                                        onChange={e => {
                                            if (e.target.checked) {
                                                setSelectedRegistrations([...selectedRegistrations, registration._id]);
                                            } else {
                                                setSelectedRegistrations(
                                                    selectedRegistrations.filter(id => id !== registration._id)
                                                );
                                            }
                                        }}
                                    />
                                </td>
                                <td className="admin_btns">
                                    {registration.status === "Ожидание" && (
                                        <>
                                            <button
                                                className="approve-btn"
                                                onClick={() => handleApprove(registration._id)}
                                            >
                                                Одобрить
                                            </button>
                                            <button
                                                className="reject-btn"
                                                onClick={() => handleReject(registration._id)}
                                            >
                                                Отклонить
                                            </button>
                                            <button
                                                className="reject-btn"
                                                onClick={() => handleDelete(registration._id)}
                                            >
                                                Удалить
                                            </button>
                                        </>
                                    )}
                                    {registration.status === "Отклоненный" && (
                                        <button
                                            className="reject-btn"
                                            onClick={() => handleDelete(registration._id)}
                                        >
                                            Удалить
                                        </button>
                                    )}
                                    {registration.status === "Одобренный" && registration.contractFilePath && (
                                        <>
                                            <button
                                                className="reject-btn"
                                                onClick={() => handleReject(registration._id)}
                                            >
                                                Отклонить
                                            </button>
                                            <button
                                                className="approve-btn"
                                                onClick={() => handleViewPdf(registration._id)}
                                            >
                                                Проверить договор
                                            </button>
                                            <button
                                                className="reject-btn"
                                                onClick={() => handleDelete(registration._id)}
                                            >
                                                Удалить
                                            </button>
                                        </>
                                    )}
                                    {registration.status === "Одобренный" && !registration.contractFilePath && (
                                        <>
                                            <button
                                                className="reject-btn"
                                                onClick={() => handleReject(registration._id)}
                                            >
                                                Отклонить
                                            </button>
                                            <button
                                                className="reject-btn"
                                                onClick={() => handleDelete(registration._id)}
                                            >
                                                Удалить
                                            </button>
                                        </>
                                    )}
                                    {registration.status === "Принят" && (
                                        <button
                                            className="reject-btn"
                                            onClick={() => handleReject(registration._id)}
                                        >
                                            Отклонить
                                        </button>
                                    )}
                                    {registration.status !== "Ожидание" &&
                                        registration.status !== "Отклоненный" &&
                                        registration.status !== "Одобренный" && (
                                            <button
                                                onClick={() => handleViewConsent(registration.userId)}
                                                className="toggle-info-button"
                                            >
                                                Просмотр согласия
                                            </button>
                                        )}
                                </td>
                            </tr>
                        );
                    })}
                </tbody>
            </table>
            {selectedRejectRegistrationId && (
                <div className="modal">
                    <div className="modal-content">
                        <h3>Причина отклонения</h3>
                        <textarea
                            value={rejectReason}
                            onChange={e => setRejectReason(e.target.value)}
                            placeholder="Причина отклонения"
                        />
                        <div className="form-buttons">
                            <button className="approve-btn" onClick={confirmReject}>
                                Подтвердить
                            </button>
                            <button
                                className="reject-btn"
                                onClick={() => {
                                    setSelectedRejectRegistrationId(null);
                                    setRejectReason("");
                                }}
                            >
                                Отмена
                            </button>
                        </div>
                    </div>
                </div>
            )}
            {selectedPdfRegistrationId && (
                <div className="modal">
                    <div className="modal-content">
                        <h3>Проверка договора</h3>
                        <iframe
                            src={pdfUrl}
                            width="100%"
                            height="500px"
                            style={{ border: "none" }}
                            title="Contract PDF"
                        />
                        <div className="form-buttons">
                            <button
                                className="approve-btn"
                                onClick={() => handleApprovePdf(selectedPdfRegistrationId)}
                            >
                                Принять
                            </button>
                            <button
                                className="reject-btn"
                                onClick={() => {
                                    setSelectedPdfRegistrationId(null);
                                    setPdfUrl("");
                                }}
                            >
                                Закрыть
                            </button>
                        </div>
                    </div>
                </div>
            )}
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
                        <div className="form-buttons">
                            <button className="reject-btn" onClick={handleCloseUserInfo}>
                                Закрыть
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AdminCoursesManagement;