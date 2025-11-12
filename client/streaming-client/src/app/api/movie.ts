import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { IMovie } from '../interfaces/movie';
import { environment } from '../../environments/environment';

@Injectable({ providedIn: 'root' })
export class MoviesApi {
  private http = inject(HttpClient);
  private apiURL = environment.apiURL;

  getMovies() {
    return this.http.get<IMovie[]>(`${this.apiURL}/movies`);
  }
}
