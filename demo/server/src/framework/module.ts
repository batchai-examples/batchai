import { Module, Provider } from '@nestjs/common';
import { User, UserFacade, UserService } from './user';
import { RequestKontextInterceptor } from './kontext';
import { JwtAuthGuard, JwtStrategy, LocalStrategy } from './auth';
import { GithubAuthRest, GithubStrategy } from './auth.github';
import { UserRest, AuthRest } from './user';
import { APP_GUARD, APP_INTERCEPTOR } from '@nestjs/core';
import { JwtService } from '@nestjs/jwt';
import { TypeOrmModule } from '@nestjs/typeorm';
import { TerminusModule } from '@nestjs/terminus';
import { HttpModule } from '@nestjs/axios';
import { HealthRest } from './health';

export const requestKontextInterceptor = {
	provide: APP_INTERCEPTOR,
	useClass: RequestKontextInterceptor,
};

export const jwtAuthGuard = {
	provide: APP_GUARD,
	useClass: JwtAuthGuard,
};

export const entities = [User];

const providers: Provider[] = [JwtService, UserService, UserFacade, LocalStrategy, JwtStrategy, GithubStrategy];

@Module({
	imports: [TypeOrmModule.forFeature(entities), TerminusModule, HttpModule],
	controllers: [UserRest, AuthRest, GithubAuthRest, HealthRest],
	providers: providers,
	exports: providers,
})
export class FrameworkModule {}
