import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import "./admin_group_management.css";

const AdminGroupManagement = () => {
    const navigate = useNavigate();
    const [groups, setGroups] = useState([]);
    const [filteredGroups, setFilteredGroups] = useState([]);
    const [searchTerm, setSearchTerm] = useState("");
    const [selectedGroup, setSelectedGroup] = useState(null);
    const [groupName, setGroupName] = useState("");
    const [courseId, setCourseId] = useState("");
    const [error, setError] = useState("");
    const [successMessage, setSuccessMessage] = useState("");
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [groupMembers, setGroupMembers] = useState([]);
    const [selectedGroupId, setSelectedGroupId] = useState(null);
    const [courses, setCourses] = useState([]);
    const [groupsWithMembers, setGroupsWithMembers] = useState({});
    const [orders, setOrders] = useState([]);
    const [selectedOrderId, setSelectedOrderId] = useState(null);
    const [isEnrollModalOpen, setIsEnrollModalOpen] = useState(false);
    const [isExpelModalOpen, setIsExpelModalOpen] = useState(false);
    const [rejectReason, setRejectReason] = useState("");
    const [isRejectModalOpen, setIsRejectModalOpen] = useState(false);
    const [selectedMembers, setSelectedMembers] = useState([]);
    const [isSelectMembersModalOpen, setIsSelectMembersModalOpen] = useState(false);

    const fetchGroups = async () => {
        try {
            const response = await fetch("http://localhost:5000/groups", { credentials: "include" });
            if (!response.ok) throw new Error("Ошибка при загрузке групп");
            const data = await response.json();
            if (!data?.groups) throw new Error("Данные не получены");

            const groupsWithMembersInfo = {};
            for (const group of data.groups) {
                const membersResponse = await fetch(`http://localhost:5000/admin/group-members/${group._id}`, { credentials: "include" });
                if (!membersResponse.ok) throw new Error("Ошибка при загрузке участников группы");
                const membersData = await membersResponse.json();
                groupsWithMembersInfo[group._id] = {
                    hasMembers: membersData.members?.length > 0,
                    hasPaidMembers: membersData.members?.some(m => m.status === "Оплаченный"),
                    hasEnrolledMembers: membersData.members?.some(m => m.status === "Проходит курс"),
                    hasCompletedMembers: membersData.members?.some(m => m.status === "Завершил"),
                    hasRejectedMembers: membersData.members?.some(m => m.status === "Отклоненный"),
                    status: group.status,
                    currentStudents: group.currentStudents,
                    members: membersData.members || []
                };
            }

            setGroups(data.groups);
            setFilteredGroups(data.groups);
            setGroupsWithMembers(groupsWithMembersInfo);
        } catch (error) {
            console.error("Ошибка загрузки групп:", error);
            setError("Ошибка при загрузке групп");
        }
    };

    const fetchCourses = async () => {
        try {
            const response = await fetch("http://localhost:5000/courses", { credentials: "include" });
            if (!response.ok) throw new Error("Ошибка при загрузке курсов");
            const data = await response.json();
            setCourses(data || []);
        } catch (error) {
            console.error("Ошибка загрузки курсов:", error);
            setError("Ошибка при загрузке курсов");
        }
    };

    const fetchOrders = async () => {
        try {
            const response = await fetch("http://localhost:5000/admin/orders", { credentials: "include" });
            if (!response.ok) throw new Error("Ошибка при загрузке приказов");
            const data = await response.json();
            setOrders(data || []);
        } catch (error) {
            console.error("Ошибка загрузки приказов:", error);
            setError("Ошибка при загрузке приказов");
        }
    };

    useEffect(() => {
        fetchGroups();
        fetchCourses();
        fetchOrders();
    }, []);

    const handleSearchChange = (e) => {
        const term = e.target.value.toLowerCase();
        setSearchTerm(term);
        setFilteredGroups(groups.filter(group =>
            group.groupName.toLowerCase().includes(term)
        ));
    };

    const handleEdit = (group) => {
        setSelectedGroup(group);
        setGroupName(group.groupName);
        setCourseId(group.courseId);
        setIsModalOpen(true);
    };

    const handleSave = async () => {
        if (!groupName || !courseId) {
            setError("Заполните название группы и выберите курс");
            return;
        }
        try {
            const response = await fetch(`http://localhost:5000/admin/update-group/${selectedGroup._id}`, {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ groupName, courseId })
            });
            const data = await response.json();
            if (data.success) {
                setSuccessMessage("Группа обновлена!");
                handleCloseAllModals();
                fetchGroups();
            } else {
                setError(data.error || "Ошибка при обновлении группы");
            }
        } catch (error) {
            console.error("Ошибка обновления группы:", error);
            setError("Ошибка сервера");
        }
    };

    const handleDelete = async (groupId) => {
        try {
            const response = await fetch(`http://localhost:5000/admin/delete-group/${groupId}`, {
                method: "DELETE",
                credentials: "include"
            });
            const data = await response.json();
            if (data.success) {
                setSuccessMessage("Группа удалена!");
                fetchGroups();
            } else {
                setError(data.error || "Ошибка при удалении группы");
            }
        } catch (error) {
            console.error("Ошибка удаления группы:", error);
            setError("Ошибка сервера");
        }
    };

    const handleRejectGroup = (groupId) => {
        setSelectedGroupId(groupId);
        setIsRejectModalOpen(true);
    };

    const confirmReject = async () => {
        if (!rejectReason) {
            setError("Укажите причину отклонения");
            return;
        }
        try {
            const response = await fetch(`http://localhost:5000/admin/reject-group-registrations/${selectedGroupId}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ reason: rejectReason })
            });
            const data = await response.json();
            if (data.success) {
                setSuccessMessage("Записи группы отклонены!");
                handleCloseAllModals();
                fetchGroups();
            } else {
                setError(data.error || "Ошибка при отклонении записей");
            }
        } catch (error) {
            console.error("Ошибка отклонения записей:", error);
            setError("Ошибка сервера");
        }
    };

    const handleCreateGroup = async (e) => {
        e.preventDefault();
        if (!groupName || !courseId) {
            setError("Заполните название группы и выберите курс");
            return;
        }
        try {
            const response = await fetch("http://localhost:5000/admin/create-group", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ groupName, courseId })
            });
            const data = await response.json();
            if (data.success) {
                setSuccessMessage("Группа создана!");
                handleCloseAllModals();
                fetchGroups();
            } else {
                setError(data.error || "Ошибка создания группы");
            }
        } catch (error) {
            console.error("Ошибка создания группы:", error);
            setError("Ошибка сервера");
        }
    };

    const fetchGroupMembers = async (groupId) => {
        try {
            const response = await fetch(`http://localhost:5000/admin/group-members/${groupId}`, { credentials: "include" });
            if (!response.ok) throw new Error("Ошибка при загрузке участников");
            const data = await response.json();
            setGroupMembers(data.members || []);
        } catch (error) {
            console.error("Ошибка загрузки участников:", error);
            setError("Ошибка загрузки участников");
        }
    };

    const handleViewMembers = (groupId) => {
        setSelectedGroupId(groupId);
        setSelectedMembers([]);
        fetchGroupMembers(groupId);
    };

    const handleExpelMember = (registrationId) => {
        setSelectedGroupId(`registration-${registrationId}`);
        setIsExpelModalOpen(true);
    };

    const handleEnrollGroup = (groupId) => {
        setSelectedGroupId(groupId);
        setIsEnrollModalOpen(true);
    };

    const handleExpelGroup = (groupId, isSelective = false) => {
        setSelectedGroupId(groupId);
        if (isSelective) {
            setIsSelectMembersModalOpen(true);
        } else {
            setIsExpelModalOpen(true);
        }
    };

    const handleSelectMembersModal = () => {
        setIsSelectMembersModalOpen(true);
    };

    const confirmEnroll = async () => {
        if (!selectedOrderId) {
            setError("Выберите приказ о зачислении");
            return;
        }
        try {
            const response = await fetch(`http://localhost:5000/admin/enroll-group/${selectedGroupId}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ orderId: selectedOrderId })
            });
            const data = await response.json();
            if (data.success) {
                setSuccessMessage("Группа зачислена!");
                handleCloseAllModals();
                fetchGroups();
            } else {
                setError(data.error || "Ошибка при зачислении группы");
            }
        } catch (error) {
            console.error("Ошибка зачисления группы:", error);
            setError("Ошибка сервера");
        }
    };

    const confirmExpel = async () => {
        if (!selectedOrderId) {
            setError("Выберите приказ");
            return;
        }
        try {
            let url, body;
            if (selectedMembers.length > 0) {
                url = `http://localhost:5000/admin/expel-registrations`;
                body = JSON.stringify({ registrationIds: selectedMembers, orderId: selectedOrderId });
            } else if (selectedGroupId.includes("registration")) {
                const registrationId = selectedGroupId.split("-")[1];
                url = `http://localhost:5000/admin/expel-registration/${registrationId}`;
                body = JSON.stringify({ orderId: selectedOrderId });
            } else {
                url = `http://localhost:5000/admin/expel-group/${selectedGroupId}`;
                body = JSON.stringify({ orderId: selectedOrderId });
            }

            const response = await fetch(url, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body
            });
            const data = await response.json();
            if (data.success) {
                setSuccessMessage(
                    selectedMembers.length > 0
                        ? `Выбранные участники (${data.count}) отчислены!`
                        : selectedGroupId.includes("registration")
                            ? "Участник отчислен!"
                            : "Группа завершена!"
                );
                handleCloseAllModals();
                fetchGroups();
            } else {
                setError(data.error || "Ошибка при отчислении");
            }
        } catch (error) {
            console.error("Ошибка отчисления:", error);
            setError("Ошибка сервера");
        }
    };

    const handleSelectMember = (memberId) => {
        setSelectedMembers(prev =>
            prev.includes(memberId)
                ? prev.filter(id => id !== memberId)
                : [...prev, memberId]
        );
    };

    const confirmSelectMembers = () => {
        if (selectedMembers.length === 0) {
            setError("Выберите хотя бы одного участника");
            return;
        }
        setIsSelectMembersModalOpen(false);
        setIsExpelModalOpen(true);
    };

    const handleCloseAllModals = () => {
        setIsModalOpen(false);
        setIsEnrollModalOpen(false);
        setIsExpelModalOpen(false);
        setIsRejectModalOpen(false);
        setIsSelectMembersModalOpen(false);
        setSelectedGroup(null);
        setSelectedGroupId(null);
        setGroupName("");
        setCourseId("");
        setGroupMembers([]);
        setSelectedMembers([]);
        setSelectedOrderId(null);
        setRejectReason("");
        setError("");
        setSuccessMessage("");
    };

    useEffect(() => {
        const handleKeyDown = e => {
            if (e.key === "Escape") handleCloseAllModals();
        };
        window.addEventListener("keydown", handleKeyDown);
        return () => window.removeEventListener("keydown", handleKeyDown);
    }, []);

    return (
        <div className="admin-groups-management-page">
            {error && <div className="error-message">{error}</div>}
            {successMessage && <div className="success-message">{successMessage}</div>}

            <div className="admin-course-header">
                <input
                    type="text"
                    placeholder="Поиск по группам..."
                    value={searchTerm}
                    onChange={handleSearchChange}
                    className="search-input"
                />
                <button className="approve-btn" onClick={() => setIsModalOpen(true)}>
                    Создать группу
                </button>
            </div>

            {isModalOpen && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>{selectedGroup ? "Редактирование группы" : "Создание группы"}</h2>
                        <form onSubmit={selectedGroup ? handleSave : handleCreateGroup}>
                            <label>Название группы:</label>
                            <input
                                type="text"
                                value={groupName}
                                onChange={e => setGroupName(e.target.value)}
                                placeholder="Введите название группы"
                            />
                            <label>Курс:</label>
                            <select value={courseId} onChange={e => setCourseId(e.target.value)}>
                                <option value="">Выберите курс</option>
                                {courses.map(course => (
                                    <option key={course._id} value={course._id}>
                                        {course.title} (Макс. студентов: {course.maxStudents})
                                    </option>
                                ))}
                            </select>
                            <div className="form-buttons">
                                <button className="approve-btn" type="submit">
                                    {selectedGroup ? "Сохранить" : "Создать"}
                                </button>
                                <button className="reject-btn" type="button" onClick={handleCloseAllModals}>
                                    Закрыть
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            <table>
                <thead>
                    <tr>
                        <th>Название группы</th>
                        <th>Курс</th>
                        <th>Студенты</th>
                        <th>Статус</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    {filteredGroups.map(group => {
                        const groupInfo = groupsWithMembers[group._id] || {};
                        const isCompleted = groupInfo.status === "Завершена";
                        const hasActiveMembers = groupInfo.hasEnrolledMembers || groupInfo.hasCompletedMembers;
                        const hasOnlyPaid = groupInfo.hasPaidMembers && !hasActiveMembers;
                        const isEmpty = !groupInfo.hasMembers || (groupInfo.hasRejectedMembers && !hasActiveMembers && !groupInfo.hasPaidMembers);

                        return (
                            <tr key={group._id}>
                                <td>{group.groupName}</td>
                                <td>{courses.find(c => c._id === group.courseId)?.title || "Неизвестный курс"}</td>
                                <td>
                                    {`${groupInfo.currentStudents || 0} / ${courses.find((course) => course._id === group.courseId)?.maxStudents || "Не указано"}`}
                                </td>
                                <td>{groupInfo.status || "Активна"}</td>
                                <td>
                                    <div className="bt">
                                        <button className="approve-btn" onClick={() => handleViewMembers(group._id)}>
                                            Просмотр участников
                                        </button>
                                        {!isCompleted && (
                                            <>
                                                {isEmpty ? (
                                                    <>
                                                        <button className="approve-btn" onClick={() => handleEdit(group)}>
                                                            Редактировать
                                                        </button>
                                                        <button className="reject-btn" onClick={() => handleRejectGroup(group._id)}>
                                                            Отклонить записи
                                                        </button>
                                                        <button className="reject-btn" onClick={() => handleDelete(group._id)}>
                                                            Удалить
                                                        </button>
                                                    </>
                                                ) : hasOnlyPaid ? (
                                                    <>
                                                        <button className="approve-btn" onClick={() => handleEnrollGroup(group._id)}>
                                                            Зачислить группу
                                                        </button>
                                                        <button className="reject-btn" onClick={() => handleRejectGroup(group._id)}>
                                                            Отклонить записи
                                                        </button>
                                                        {!groupInfo.hasEnrolledMembers && (
                                                            <button className="reject-btn" onClick={() => handleDelete(group._id)}>
                                                                Удалить
                                                            </button>
                                                        )}
                                                    </>
                                                ) : hasActiveMembers ? (
                                                    <button className="approve-btn" onClick={() => handleExpelGroup(group._id, false)}>
                                                        Завершить
                                                    </button>
                                                ) : null}
                                            </>
                                        )}
                                    </div>
                                </td>
                            </tr>
                        );
                    })}
                </tbody>
            </table>

            {selectedGroupId && !selectedGroupId.includes("registration") && !isSelectMembersModalOpen && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Участники группы</h2>
                        <table>
                            <thead>
                                <tr>
                                    <th>Имя пользователя</th>
                                    <th>Email</th>
                                    <th>Статус</th>
                                    <th>Действия</th>
                                </tr>
                            </thead>
                            <tbody>
                                {groupMembers.map(member => (
                                    <tr key={member._id}>
                                        <td>{member.username}</td>
                                        <td>{member.email}</td>
                                        <td>{member.status || "Неизвестно"}</td>
                                        <td>
                                            <div className="bt">
                                                <button
                                                    className="toggle-info-button"
                                                    onClick={() => navigate(`/admin/profile?username=${member.username}`)}
                                                >
                                                    Данные
                                                </button>
                                                {member.status === "Проходит курс" && (
                                                    <button className="reject-btn" onClick={() => handleExpelMember(member._id)}>
                                                        Отчислить
                                                    </button>
                                                )}
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                        <div className="button">
                            {groupMembers.some(m => m.status === "Проходит курс") && (
                                <button className="reject-btn" onClick={handleSelectMembersModal}>
                                    Отчислить выбранных
                                </button>
                            )}
                            <button className="reject-btn" onClick={handleCloseAllModals}>
                                Закрыть
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {isSelectMembersModalOpen && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Выберите участников для отчисления</h2>
                        <table>
                            <thead>
                                <tr>
                                    <th>Выбрать</th>
                                    <th>Имя пользователя</th>
                                    <th>Email</th>
                                    <th>Статус</th>
                                </tr>
                            </thead>
                            <tbody>
                                {groupMembers.filter(m => m.status === "Проходит курс").map(member => (
                                    <tr key={member._id}>
                                        <td>
                                            <input
                                                type="checkbox"
                                                checked={selectedMembers.includes(member._id)}
                                                onChange={() => handleSelectMember(member._id)}
                                            />
                                        </td>
                                        <td>{member.username}</td>
                                        <td>{member.email}</td>
                                        <td>{member.status}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                        <div className="button">
                            <button className="approve-btn" onClick={confirmSelectMembers}>
                                Подтвердить выбор
                            </button>
                            <button className="reject-btn" onClick={handleCloseAllModals}>
                                Отмена
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {isEnrollModalOpen && (
                <div className="modal">
                    <div className="modal-content">
                        <h3>Зачисление группы</h3>
                        <div className="form-group">
                            <label>Приказ о зачислении:</label>
                            <select
                                value={selectedOrderId || ""}
                                onChange={(e) => setSelectedOrderId(e.target.value)}
                            >
                                <option value="">Выберите приказ</option>
                                {orders
                                    .filter((order) => order.orderType === "О зачислении обучающихся")
                                    .map((order) => (
                                        <option key={order.id} value={order.id}>
                                            {order.number} от {new Date(order.date).toLocaleDateString()} ({order.orderType})
                                        </option>
                                    ))}
                            </select>
                        </div>
                        <button className="approve-btn" onClick={confirmEnroll}>
                            Подтвердить
                        </button>
                        <button className="reject-btn" onClick={handleCloseAllModals}>
                            Отмена
                        </button>
                    </div>
                </div>
            )}

            {isExpelModalOpen && (
                <div className="modal">
                    <div className="modal-content">
                        <h3>Выбор приказа</h3>
                        <div className="form-group">
                            <label>Приказ:</label>
                            <select
                                value={selectedOrderId || ""}
                                onChange={(e) => setSelectedOrderId(e.target.value)}
                            >
                                <option value="">Выберите приказ</option>
                                {orders
                                    .filter(
                                        (order) =>
                                            order.orderType === "О выпуске обучающихся" ||
                                            order.orderType === "Об отчислении обучающихся"
                                    )
                                    .map((order) => (
                                        <option key={order.id} value={order.id}>
                                            {order.number} от {new Date(order.date).toLocaleDateString()} ({order.orderType})
                                        </option>
                                    ))}
                            </select>
                        </div>
                        <button className="approve-btn" onClick={confirmExpel}>
                            Подтвердить
                        </button>
                        <button className="reject-btn" onClick={handleCloseAllModals}>
                            Отмена
                        </button>
                    </div>
                </div>
            )}

            {isRejectModalOpen && (
                <div className="modal">
                    <div className="modal-content">
                        <h3>Укажите причину отклонения</h3>
                        <textarea
                            value={rejectReason}
                            onChange={(e) => setRejectReason(e.target.value)}
                            placeholder="Причина отклонения (например, мало людей)"
                        />
                        <div className="button">
                            <button className="approve-btn" onClick={confirmReject}>
                                Подтвердить
                            </button>
                            <button className="reject-btn" onClick={handleCloseAllModals}>
                                Отмена
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AdminGroupManagement;