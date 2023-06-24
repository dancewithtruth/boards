import { BoardWithMembers } from '@/api';
import { MEMBERSHIP_ROLES } from '@/constants';

export const getMaxFieldFromObj = <T, K extends keyof T & string>(obj: { [key: string]: T }, field: K) => {
  let maxNumber = 0;
  for (const key in obj) {
    const value = obj[key][field];
    if (typeof value === 'number' && value > maxNumber) {
      maxNumber = value;
    }
  }
  return maxNumber;
};

export const displayColor = (str: string): string => {
  return str.toLowerCase().replace(/_/g, ' ');
};

export const mergeArrays = <T extends Record<any, any>>(fieldName: keyof T, ...arrays: T[][]): T[] => {
  const mergedArray: T[] = arrays.flat();

  const mergedObject = mergedArray.reduce((result, item) => {
    const itemId = item[fieldName];
    if (!result[itemId]) {
      result[itemId] = item;
    }
    return result;
  }, {} as Record<string, T>);

  return Object.values(mergedObject);
};

export const isAdmin = (userID: string, boardWithMembers: BoardWithMembers): boolean => {
  if (boardWithMembers.user_id == userID) {
    return true;
  }
  return boardWithMembers.members.some((member) => {
    if (member.id == userID && member.membership.role == MEMBERSHIP_ROLES.ADMIN) {
      return true;
    }
  });
};
