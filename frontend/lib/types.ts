/**
 * Типы, отражающие модели Go-бэкенда (internal/model).
 */

export type QuestionType = "single" | "multiple" | "text";

export interface Question {
  id: number;
  type: QuestionType;
  title: string;
  description: string;
  image: string;
  options: string[];
}

export interface PollConfig {
  questions: Question[];
}

export interface Poll {
  id: number;
  title: string;
  description: string;
  config: PollConfig;
  creator_id: number;
  short_id: string;
  secured: boolean;
  auth_only: boolean;
  edited_at: string;
  created_at: string;
}

export interface NewPollRequest {
  title: string;
  description: string;
  config: PollConfig;
  secured: boolean;
  auth_only: boolean;
}

export interface UpdatePollRequest extends NewPollRequest {}

export interface Answer {
  question_id: number;
  options: string[];
}

export interface NewVoteRequest {
  answers: Answer[];
}

export interface UserResponse {
  id: number;
  username: string;
  created_at: string;
}

// --- Статистика ---

export interface Result {
  option: string;
  votes: number;
  percentage: number;
}

export interface QuestionResult {
  id: number;
  options: Result[];
}

export interface CountryResult {
  country_code: string;
  votes: number;
}

export interface Stats {
  total_votes: number;
  results: QuestionResult[];
  top_countries: CountryResult[];
}
