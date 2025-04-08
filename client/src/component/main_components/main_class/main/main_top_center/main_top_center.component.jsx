import React from "react";
import Slider from "react-slick";
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import "./main_top_center.component.css";
import apt from "../../../../../assets/apt.png";

const Main_top_center = () => {
    const settings = {
        dots: false,
        infinite: true,
        speed: 500,
        slidesToShow: 1,
        slidesToScroll: 1,
        autoplay: true,
        autoplaySpeed: 3000,
        arrows: false,
        pauseOnHover: true,
        centerMode: false,
        variableWidth: false,
    };

    return (
        <div className="main_top_for_slider">
            <div className="main_top_center_slider">
                {/* Фоновое изображение */}
                <div className="main_top_img">
                    <img src={apt} alt="Наш техникум" />
                </div>

                {/* Контейнер слайдера */}
                <div className="slider_container_wrapper">
                    <div className="main_top_slider">
                        <Slider {...settings}>
                            {/* Слайд 1 */}
                            <div className="slide_wrapper">
                                <div className="slider-item">
                                    <div className="slider-content">
                                        <h3>Почему наш сайт?</h3>
                                        <p>
                                            Мы предлагаем качественные образовательные программы, разработанные опытными преподавателями.
                                        </p>
                                    </div>
                                </div>
                            </div>

                            {/* Слайд 2 - специальный */}
                            <div className="slide_wrapper">
                                <div className="slider-item special-slide">
                                    <div className="slider-content">
                                        <h3>Что нужно для профессиональной переподготовке?</h3>
                                        <ul>
                                            <li>Лица, имеющие среднее профессиональное и (или) высшее образование</li>
                                            <li>Лица, получающие среднее профессиональное и (или) высшее образование</li>
                                        </ul>
                                    </div>
                                </div>
                            </div>

                            {/* Слайд 3 */}
                            <div className="slide_wrapper">
                                <div className="slider-item">
                                    <div className="slider-content">
                                        <h3>Когда и где?</h3>
                                        <p>Запись на курс 24/7</p>
                                        <p>Возможность обучения в любой момент времени, не дожидаясь формирования группы</p>
                                    </div>
                                </div>
                            </div>
                            {/* Слайд 2 - специальный */}
                            <div className="slide_wrapper">
                                <div className="slider-item special-slide">
                                    <div className="slider-content">
                                        <h3>Наше преимущество</h3>
                                        <p>Возможность обучения независимо от места проживания</p>
                                    </div>
                                </div>
                            </div>
                        </Slider>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Main_top_center;