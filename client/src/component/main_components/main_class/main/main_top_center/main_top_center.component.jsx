import React from "react";
import Slider from "react-slick";
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import "./main_top_center.component.css";
import apt from "../../../../../assets/apt.png";

const Main_top_center = () => {
    const settings = {
        dots: true,
        infinite: true,
        speed: 500,
        slidesToShow: 1,
        slidesToScroll: 1,
        autoplay: true,
        autoplaySpeed: 3000,
        arrows: true, // Добавить стрелки навигации
        pauseOnHover: true, // Пауза при наведении
        centerMode: true, // Центрирование слайдов
        centerPadding: "0", // Убрать отступы по бокам
    };

    return (
        <>

        <div className="main_top_for_slider">
            <div className="main_top_center_slider">
                {/* Картинка как фон */}
                <div className="main_top_img">
                    <img src={apt} alt="Наш техникум" />
                </div>

                {/* Слайдер с текстовым содержимым */}
                <div className="main_top_slider">
                    <Slider {...settings}>
                        {/* Слайд 1 */}
                        <div className="slider-item">
                            <div className="slider-content">
                                <h3>Почему наш сайт?</h3>
                                <p>
                                    Мы предлагаем качественные образовательные программы, разработанные опытными преподавателями.
                                </p>
                            </div>
                        </div>

                        {/* Слайд 2 */}
                        <div className="slider-item">
                            <div className="slider-content">
                                <h3>Почему именно у нас?</h3>
                                <p>
                                    По завершению обучения выдается диплом о профессиональной переподготовке.
                                </p>
                            </div>
                        </div>

                        {/* Слайд 3 */}
                        <div className="slider-item">
                            <div className="slider-content">
                                <h3>Почему именно у нас?</h3>
                                <p>
                                    Наши курсы помогут вам получить актуальные знания и навыки, необходимые для успешной карьеры.
                                </p>
                            </div>
                        </div>

                        {/* Слайд 4 */}
                        <div className="slider-item">
                            <div className="slider-content">
                                <h3>Почему наш мы?</h3>
                                <p>
                                    Доставка документов по всей России.
                                </p>
                            </div>
                        </div>

                        {/* Слайд 5 */}
                        <div className="slider-item">
                            <div className="slider-content">
                                <h3>Преимущества наших курсов</h3>
                                <p>
                                    Гибкий график обучения, доступные цены и поддержка на всех этапах обучения.
                                </p>
                            </div>
                        </div>
                    </Slider>
                </div>

                {/* Текст под слайдером */}
                
            </div>
           
        </div>
         <div className="under_img">
                    <p>
                        К освоению дополнительных профессиональных образовательных программ допускаются:
                    </p>
                    <p>
                        1. лица, имеющие среднее профессиональное и (или) высшее образование;
                    </p>
                    <p>
                        2. лица, получающие среднее профессиональное и (или) высшее образование.
                    </p>
                </div>
        </>
    );
};

export default Main_top_center;