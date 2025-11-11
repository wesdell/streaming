export interface IMovieRanking {
  ranking_name: string;
}

export interface IMovie {
  _id: string;
  imdb_id: string;
  youtube_id: string;
  title: string;
  poster_path: string;
  ranking?: IMovieRanking;
  review?: string;
}
