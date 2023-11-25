import useSWR from "swr";
import { Box, List, ThemeIcon, MantineProvider } from '@mantine/core'
// import './App.css';
import '@mantine/core/styles.css';
import AddTodo from './components/AddTodo';
import { CheckCircleFillIcon } from "@primer/octicons-react";

export interface Todo {
    id: number;
    title: string;
    body: string;
    done: boolean;
}

export const ENDPOINT = 'http://localhost:8080'

const fetcher = (url: string) => fetch(`${ENDPOINT}/${url}`).then(r => r.json());

function App() {

    const { data, mutate } = useSWR<Todo[]>('api/todos', fetcher, { revalidateOnMount: true });

    async function markTodoAsDone(id: number) {
        const updated = await fetch(`${ENDPOINT}/api/todos/${id}/done`, {
            method: "PATCH",
        }).then((r) => r.json());

        mutate(updated);
    }

    return <MantineProvider>{
        <Box m="xl" w="xl" p="xl">
            <List spacing="xl" size="xl" mb={12} center>
                {data?.map((todo) => {
                    return (
                        <List.Item maw="xl"
                            onClick={() => markTodoAsDone(todo.id)}
                            key={`todo_list__${todo.id}`}
                            icon={
                                todo.done ? (
                                    <ThemeIcon color="teal" size={24} radius="xl">
                                        <CheckCircleFillIcon size={20} />
                                    </ThemeIcon>
                                ) : (
                                    <ThemeIcon color="grey" size={24} radius="xl">
                                        <CheckCircleFillIcon size={20} />
                                    </ThemeIcon>
                                )
                            }
                        >
                            {todo.title}
                        </List.Item>
                    );
                })}
            </List>
            <AddTodo mutate={mutate} />
        </Box>
    }</MantineProvider>;
}

export default App;