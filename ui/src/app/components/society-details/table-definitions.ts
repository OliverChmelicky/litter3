export const societyRequestsColumnDefinition: string[] = [
  'email',
  'first-name',
  'last-name',
  'button',
]

export const EventsColumnsDefinition: string[] = [
  'date',
  'interested-people',
  'details',
]

export const membersColumnsDefinition: string[] = [
  'image',
  'first-name',
  'last-name',
]

export const membersColumnsSocietyEditDefinition: string[] = [
  'image',
  'first-name',
  'last-name',
  'email',
  'roles',
  'deleteUser',
]

interface rolesInterface {
  key:string,
  value: string,
}

export const roles: rolesInterface[] = [
  {
    key:'admin',
    value: 'Administrator',
  },
  {
    key:'editor',
    value: 'Editor',
  },
  {
    key:'member',
    value: 'Member',
  },
]
