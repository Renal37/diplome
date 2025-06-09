import React, { useState, useEffect } from 'react';
import './admin_order.css';

const AdminOrders = () => {
    const [orderTypes, setOrderTypes] = useState([]);
    const [orders, setOrders] = useState([]);
    const [error, setError] = useState('');
    const [newOrder, setNewOrder] = useState({
        number: '',
        date: new Date().toISOString().slice(0, 10),
        orderTypeId: ''
    });
    const [newOrderType, setNewOrderType] = useState('');
    const [editingOrderType, setEditingOrderType] = useState(null);
    const [editingOrder, setEditingOrder] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = async () => {
        try {
            const [typesRes, ordersRes] = await Promise.all([
                fetch('http://localhost:5000/admin/order-types', { credentials: 'include' }),
                fetch('http://localhost:5000/admin/orders', { credentials: 'include' })
            ]);

            const typesData = await typesRes.json();
            const ordersData = await ordersRes.json();

            setOrderTypes(typesData || []);
            setOrders(ordersData || []);
        } catch (error) {
            console.error('Ошибка получения данных:', error);
            setError('Не удалось загрузить данные');
        } finally {
            setLoading(false);
        }
    };

    const handleAddOrderType = async () => {
        if (!newOrderType.trim()) {
            setError('Введите название типа приказа');
            return;
        }

        try {
            const response = await fetch('http://localhost:5000/admin/order-types/add', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ name: newOrderType }),
            });

            const data = await response.json();
            if (!response.ok) {
                throw new Error(data.error || 'Ошибка сервера');
            }

            setOrderTypes(prev => [...prev, data]);
            setNewOrderType('');
            setError('');
        } catch (error) {
            console.error('Ошибка добавления типа приказа:', error);
            setError(error.message);
        }
    };

    const isValidDate = (dateString) => {
        const regEx = /^\d{4}-\d{2}-\d{2}$/;
        if (!dateString.match(regEx)) return false;
        const d = new Date(dateString);
        return d instanceof Date && !isNaN(d);
    };

    const handleAddOrder = async () => {
        if (!isValidDate(newOrder.date)) {
            setError('Укажите корректную дату в формате ГГГГ-ММ-ДД');
            return;
        }

        if (!newOrder.orderTypeId) {
            setError('Выберите тип приказа');
            return;
        }

        try {
            const response = await fetch('http://localhost:5000/admin/orders/add', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({
                    number: newOrder.number,
                    date: newOrder.date,
                    orderTypeId: newOrder.orderTypeId
                }),
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Ошибка сервера');
            }

            setOrders(prev => [...prev, data]);
            setNewOrder({
                number: '',
                date: new Date().toISOString().slice(0, 10),
                orderTypeId: ''
            });
            setError('');
        } catch (error) {
            console.error('Ошибка добавления приказа:', error);
            setError(error.message);
        }
    };

    const handleUpdateOrderType = async () => {
        if (!editingOrderType || !editingOrderType.name.trim()) {
            setError('Введите название типа приказа');
            return;
        }

        try {
            const response = await fetch(`http://localhost:5000/admin/order-types/update?id=${editingOrderType.id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ name: editingOrderType.name }),
            });

            const data = await response.json();
            if (!response.ok) {
                throw new Error(data.error || 'Ошибка сервера');
            }

            setOrderTypes(prev => prev.map(type =>
                type.id === editingOrderType.id ? data : type
            ));
            setEditingOrderType(null);
            setError('');
        } catch (error) {
            console.error('Ошибка обновления типа приказа:', error);
            setError(error.message);
        }
    };

    const handleDeleteOrderType = async (id) => {
        if (!window.confirm('Удалить этот тип приказа и все связанные приказы?')) return;

        try {
            const response = await fetch(`http://localhost:5000/admin/order-types/delete?id=${id}`, {
                method: 'DELETE',
                credentials: 'include',
            });

            if (!response.ok) {
                throw new Error('Ошибка при удалении');
            }

            setOrderTypes(prev => prev.filter(type => type.id !== id));
            setOrders(prev => prev.filter(order => order.orderTypeId !== id));
            setError('');
        } catch (error) {
            console.error('Ошибка удаления типа приказа:', error);
            setError(error.message);
        }
    };

    const handleUpdateOrder = async () => {
        if (!editingOrder) return;

        if (!isValidDate(editingOrder.date)) {
            setError('Укажите корректную дату в формате ГГГГ-ММ-ДД');
            return;
        }

        if (!editingOrder.orderTypeId) {
            setError('Выберите тип приказа');
            return;
        }

        try {
            const response = await fetch(`http://localhost:5000/admin/orders/update?id=${editingOrder.id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({
                    number: editingOrder.number,
                    date: editingOrder.date,
                    orderTypeId: editingOrder.orderTypeId
                }),
            });

            const data = await response.json();
            if (!response.ok) {
                throw new Error(data.error || 'Ошибка сервера');
            }

            setOrders(prev => prev.map(order =>
                order.id === editingOrder.id ? data : order
            ));
            setEditingOrder(null);
            setError('');
        } catch (error) {
            console.error('Ошибка обновления приказа:', error);
            setError(error.message);
        }
    };

    const handleDeleteOrder = async (id) => {
        if (!window.confirm('Вы уверены, что хотите удалить этот приказ?')) return;

        try {
            const response = await fetch(`http://localhost:5000/admin/orders/delete?id=${id}`, {
                method: 'DELETE',
                credentials: 'include',
            });

            const data = await response.json();
            if (!response.ok) {
                throw new Error(data.error || 'Ошибка при удалении');
            }

            setOrders(prev => prev.filter(order => order.id !== id));
            setError('');
        } catch (error) {
            console.error('Ошибка удаления приказа:', error);
            setError(error.message);
        }
    };

    if (loading) return <div>Загрузка...</div>;

    return (
        <div className="admin-orders-container">
            {error && <div className="error-message">{error}</div>}
            <h2>Управление приказами</h2>

            <div className="order-types-section">
                <h3>Типы приказов</h3>
                <div className="add-order-type">
                    <input
                        type="text"
                        value={editingOrderType ? editingOrderType.name : newOrderType}
                        onChange={(e) => editingOrderType
                            ? setEditingOrderType({ ...editingOrderType, name: e.target.value })
                            : setNewOrderType(e.target.value)}
                        placeholder={editingOrderType ? "Редактировать тип" : "Новый тип приказа"}
                    />
                    {editingOrderType ? (
                        <>
                            <button onClick={handleUpdateOrderType}>Сохранить</button>
                            <button onClick={() => setEditingOrderType(null)}>Отмена</button>
                        </>
                    ) : (
                        <button onClick={handleAddOrderType}>Добавить</button>
                    )}
                </div>
                <ul>
                    {orderTypes.map(type => (
                        <li key={type.id}>
                            {type.name}
                            <div className="btn">

                                <button onClick={() => setEditingOrderType(type)}>Изменить</button>
                                <button onClick={() => handleDeleteOrderType(type.id)}>Удалить</button>
                            </div>

                        </li>
                    ))}

                </ul>
            </div>

            <div className="orders-section">
                <h3>{editingOrder ? "Редактирование приказа" : "Создание приказа"}</h3>
                <div className="order-form">
                    <input
                        type="text"
                        value={editingOrder ? editingOrder.number : newOrder.number}
                        onChange={(e) => editingOrder
                            ? setEditingOrder({ ...editingOrder, number: e.target.value })
                            : setNewOrder({ ...newOrder, number: e.target.value })}
                        placeholder="Номер приказа"
                    />
                    <input
                        type="date"
                        value={editingOrder ? editingOrder.date : newOrder.date}
                        onChange={(e) => editingOrder
                            ? setEditingOrder({ ...editingOrder, date: e.target.value })
                            : setNewOrder({ ...newOrder, date: e.target.value })}
                    />
                    <select
                        value={editingOrder ? editingOrder.orderTypeId : newOrder.orderTypeId}
                        onChange={(e) => editingOrder
                            ? setEditingOrder({ ...editingOrder, orderTypeId: e.target.value })
                            : setNewOrder({ ...newOrder, orderTypeId: e.target.value })}
                    >
                        <option value="">Выберите тип приказа</option>
                        {orderTypes.map(type => (
                            <option key={type.id} value={type.id}>{type.name}</option>
                        ))}
                    </select>
                    {editingOrder ? (
                        <>
                            <button onClick={handleUpdateOrder}>Сохранить</button>
                            <button onClick={() => setEditingOrder(null)}>Отмена</button>
                        </>
                    ) : (
                        <button onClick={handleAddOrder}>Создать приказ</button>
                    )}
                </div>

                <h3>Список приказов</h3>
                <table>
                    <thead>
                        <tr>
                            <th>Номер</th>
                            <th>Дата</th>
                            <th>Тип</th>
                            <th>Действия</th>
                        </tr>
                    </thead>
                    <tbody>
                        {orders.map(order => (
                            <tr key={order.id}>
                                <td>{order.number}</td>
                                <td>{new Date(order.date).toLocaleDateString()}</td>
                                <td>{order.orderType}</td>
                                <td>
                                    <button onClick={() => setEditingOrder({
                                        id: order.id,
                                        number: order.number,
                                        date: new Date(order.date).toISOString().slice(0, 10),
                                        orderTypeId: order.orderTypeId
                                    })}>Изменить</button>
                                    <button onClick={() => handleDeleteOrder(order.id)}>Удалить</button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default AdminOrders;