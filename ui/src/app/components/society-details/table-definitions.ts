export const societyRequestsColumnDefinition: string[] = [
  'email',
  'first-name',
  'last-name',
  'button',
]

export const upcommingEventsColumnsDefinition: string[] = [
  'image',
  'date',
  'interestedPeople',
]

export const participatedEventsColumnsDefinition: string[] = [
  'image',
  'date',
  'attendedPeople',
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
