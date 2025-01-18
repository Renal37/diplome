import React, { useState, useEffect } from 'react';
import { Card, CardContent } from './ui/card';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Select, SelectItem } from './ui/select';
import { Textarea } from './ui/textarea';
import { useToast } from './ui/use-toast';
import { Loader2 } from 'lucide-react';

const initialCourse = {
  name: '',
  description: '',
  price: 0,
  type: 'professional_development'
};

const AdminCourseManagement = () => {
  const [courses, setCourses] = useState([]); // Initialize as an empty array
  const [newCourse, setNewCourse] = useState(initialCourse);
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { toast } = useToast();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      const response = await fetch('http://localhost:8080/api/courses', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newCourse),
      });

      if (!response.ok) {
        throw new Error('Failed to add course');
      }

      setNewCourse(initialCourse);
      await fetchCourses();

      toast({
        title: "Success",
        description: "Course has been added successfully",
      });
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to add course. Please try again.",
        variant: "destructive",
      });
      console.error('Error adding course:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const fetchCourses = async () => {
    setIsLoading(true);
    try {
      const response = await fetch('http://localhost:8080/api/courses');
      if (!response.ok) {
        throw new Error('Failed to fetch courses');
      }
      const data = await response.json();
      setCourses(data);
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to fetch courses. Please refresh the page.",
        variant: "destructive",
      });
      console.error('Error fetching courses:', error);
      setCourses([]); // Ensure courses is set to an empty array on error
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchCourses();
  }, []);

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">Управление курсами</h1>

      <Card className="mb-6">
        <CardContent className="pt-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block mb-2">Название курса</label>
              <Input
                value={newCourse.name}
                onChange={(e) => setNewCourse({ ...newCourse, name: e.target.value })}
                required
                disabled={isSubmitting}
              />
            </div>

            <div>
              <label className="block mb-2">Описание</label>
              <Textarea
                value={newCourse.description}
                onChange={(e) => setNewCourse({ ...newCourse, description: e.target.value })}
                required
                disabled={isSubmitting}
              />
            </div>

            <div>
              <label className="block mb-2">Стоимость</label>
              <Input
                type="number"
                min="0"
                step="100"
                value={newCourse.price}
                onChange={(e) => setNewCourse({ ...newCourse, price: Number(e.target.value) })}
                required
                disabled={isSubmitting}
              />
            </div>

            <div>
              <label className="block mb-2">Тип курса</label>
              <Select
                value={newCourse.type}
                onChange={(e) => setNewCourse({ ...newCourse, type: e.target.value })}
                disabled={isSubmitting}
              >
                <SelectItem value="professional_retraining">
                  Профессиональная переподготовка
                </SelectItem>
                <SelectItem value="professional_development">
                  Повышение квалификации
                </SelectItem>
              </Select>
            </div>

            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {isSubmitting ? 'Добавление...' : 'Добавить курс'}
            </Button>
          </form>
        </CardContent>
      </Card>

      {isLoading ? (
        <div className="flex justify-center">
          <Loader2 className="h-8 w-8 animate-spin" />
        </div>
      ) : (
        <div className="grid gap-4">
          {courses && courses.length > 0 ? (
            courses.map((course, index) => (
              <Card key={index}>
                <CardContent className="pt-6">
                  <h3 className="font-bold">{course.name}</h3>
                  <p className="text-gray-600">{course.description}</p>
                  <p className="mt-2">Стоимость: {course.price.toLocaleString('ru-RU')} ₽</p>
                  <p>Тип: {
                    course.type === 'professional_retraining'
                      ? 'Профессиональная переподготовка'
                      : 'Повышение квалификации'
                  }</p>
                </CardContent>
              </Card>
            ))
          ) : (
            <p>No courses available</p>
          )}
        </div>
      )}
    </div>
  );
};

export default AdminCourseManagement;